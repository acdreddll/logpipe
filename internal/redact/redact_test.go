package redact_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/redact"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_MaskField(t *testing.T) {
	r := redact.New([]string{"password"}, "***")
	out, err := r.Apply(`{"user":"alice","password":"secret"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["password"] != "***" {
		t.Errorf("expected masked password, got %v", m["password"])
	}
	if m["user"] != "alice" {
		t.Errorf("user field altered")
	}
}

func TestApply_DeleteField(t *testing.T) {
	r := redact.New([]string{"token"}, "")
	out, err := r.Apply(`{"user":"bob","token":"abc123"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["token"]; ok {
		t.Error("token field should have been deleted")
	}
}

func TestApply_MissingField_NoError(t *testing.T) {
	r := redact.New([]string{"secret"}, "***")
	out, err := r.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("msg field altered")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r := redact.New([]string{"x"}, "***")
	_, err := r.Apply(`not json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApply_MultipleFields(t *testing.T) {
	r := redact.New([]string{"password", "ssn"}, "[REDACTED]")
	out, err := r.Apply(`{"user":"carol","password":"p4ss","ssn":"123-45-6789"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["password"] != "[REDACTED]" || m["ssn"] != "[REDACTED]" {
		t.Errorf("fields not fully redacted: %v", m)
	}
}
