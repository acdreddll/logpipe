package taginject

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyValue(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("production", WithField(""))
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_InjectsDefaultField(t *testing.T) {
	inj, err := New("production")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := inj.Apply([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["tag"] != "production" {
		t.Fatalf("expected tag=production, got %v", m["tag"])
	}
}

func TestApply_CustomField(t *testing.T) {
	inj, err := New("eu-west-1", WithField("region"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := inj.Apply([]byte(`{"level":"warn"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["region"] != "eu-west-1" {
		t.Fatalf("expected region=eu-west-1, got %v", m["region"])
	}
}

func TestApply_OverwritesExisting(t *testing.T) {
	inj, _ := New("staging")
	out, err := inj.Apply([]byte(`{"tag":"old","msg":"x"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["tag"] != "staging" {
		t.Fatalf("expected tag=staging, got %v", m["tag"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	inj, _ := New("prod")
	_, err := inj.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_PreservesExistingFields(t *testing.T) {
	inj, err := New("prod")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := inj.Apply([]byte(`{"level":"error","msg":"boom","caller":"main.go:42"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "error" {
		t.Fatalf("expected level=error, got %v", m["level"])
	}
	if m["msg"] != "boom" {
		t.Fatalf("expected msg=boom, got %v", m["msg"])
	}
	if m["caller"] != "main.go:42" {
		t.Fatalf("expected caller=main.go:42, got %v", m["caller"])
	}
	if m["tag"] != "prod" {
		t.Fatalf("expected tag=prod, got %v", m["tag"])
	}
}
