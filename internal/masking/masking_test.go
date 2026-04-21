package masking

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", `\d+`, "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("email", `[invalid`, "")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_MasksMatch(t *testing.T) {
	m, err := New("email", `[^@]+@[^@]+`, "[REDACTED]")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := m.Apply([]byte(`{"email":"user@example.com","level":"info"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	obj := decode(t, out)
	if obj["email"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", obj["email"])
	}
}

func TestApply_DefaultReplacement(t *testing.T) {
	m, err := New("token", `[A-Za-z0-9]{16,}`, "")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := m.Apply([]byte(`{"token":"abcdefghijklmnop"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	obj := decode(t, out)
	if obj["token"] != "***" {
		t.Errorf("expected ***, got %v", obj["token"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	m, err := New("secret", `.*`, "***")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	original := []byte(`{"level":"debug","msg":"ok"}`)
	out, err := m.Apply(original)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if string(out) != string(original) {
		// field absent — content may differ in key order but values must match
		obj := decode(t, out)
		if _, found := obj["secret"]; found {
			t.Error("unexpected secret field in output")
		}
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m, err := New("field", `x`, "y")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = m.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	m, err := New("count", `\d+`, "0")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := m.Apply([]byte(`{"count":42}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	obj := decode(t, out)
	if obj["count"] != float64(42) {
		t.Errorf("expected numeric 42 unchanged, got %v", obj["count"])
	}
}
