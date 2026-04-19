package enrichment

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_AddsFields(t *testing.T) {
	e, _ := New(map[string]string{"env": "prod", "region": "us-east-1"})
	out, err := e.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "prod" || m["region"] != "us-east-1" || m["msg"] != "hello" {
		t.Errorf("unexpected output: %v", m)
	}
}

func TestApply_DoesNotOverwrite(t *testing.T) {
	e, _ := New(map[string]string{"env": "prod"})
	out, err := e.Apply(`{"env":"dev"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "dev" {
		t.Errorf("expected existing value to be preserved, got %v", m["env"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	e, _ := New(map[string]string{"env": "prod"})
	_, err := e.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error when no fields provided")
	}
}
