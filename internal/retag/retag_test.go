package retag_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/retag"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := retag.New("", map[string]string{"a": "b"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyMapping(t *testing.T) {
	_, err := retag.New("level", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty mapping")
	}
}

func TestApply_RemapsValue(t *testing.T) {
	r, err := retag.New("level", map[string]string{"warn": "warning", "err": "error"})
	if err != nil {
		t.Fatal(err)
	}
	out, err := r.Apply(`{"level":"warn","msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["level"] != "warning" {
		t.Fatalf("want warning, got %v", m["level"])
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	r, _ := retag.New("level", map[string]string{"warn": "warning"})
	const in = `{"msg":"no level field"}`
	out, err := r.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Fatalf("expected unchanged line, got %s", out)
	}
}

func TestApply_NoMatchNoDefault_Unchanged(t *testing.T) {
	r, _ := retag.New("level", map[string]string{"warn": "warning"})
	const in = `{"level":"info"}`
	out, err := r.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if decode(t, out)["level"] != "info" {
		t.Fatal("expected level to remain info")
	}
}

func TestApply_NoMatchUsesDefault(t *testing.T) {
	r, _ := retag.New("level", map[string]string{"warn": "warning"}, retag.WithDefault("unknown"))
	out, err := r.Apply(`{"level":"debug"}`)
	if err != nil {
		t.Fatal(err)
	}
	if decode(t, out)["level"] != "unknown" {
		t.Fatal("expected default tag 'unknown'")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r, _ := retag.New("level", map[string]string{"a": "b"})
	_, err := r.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
