package fieldprefix_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/fieldprefix"
)

func decode(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := fieldprefix.New("", "pre_")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyPrefix(t *testing.T) {
	_, err := fieldprefix.New("env", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestApply_PrependPrefix(t *testing.T) {
	p, err := fieldprefix.New("env", "prod_")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	out, err := p.Apply([]byte(`{"env":"us-east"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	m := decode(t, out)
	if got := m["env"].(string); got != "prod_us-east" {
		t.Errorf("got %q, want %q", got, "prod_us-east")
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	p, _ := fieldprefix.New("env", "prod_")

	input := []byte(`{"level":"info"}`)
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	m := decode(t, out)
	if _, ok := m["env"]; ok {
		t.Error("env field should not be present")
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	p, _ := fieldprefix.New("count", "x")

	input := []byte(`{"count":42}`)
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	m := decode(t, out)
	if v, ok := m["count"].(float64); !ok || v != 42 {
		t.Errorf("expected count=42, got %v", m["count"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := fieldprefix.New("env", "prod_")

	_, err := p.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
