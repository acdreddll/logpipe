package tagstrip_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/tagstrip"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyFields(t *testing.T) {
	_, err := tagstrip.New([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_EmptyFieldName(t *testing.T) {
	_, err := tagstrip.New([]string{"valid", ""})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestApply_RemovesFields(t *testing.T) {
	s, err := tagstrip.New([]string{"debug", "trace"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	input := []byte(`{"msg":"hello","debug":true,"trace":"x","level":"info"}`)
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["debug"]; ok {
		t.Error("expected 'debug' to be removed")
	}
	if _, ok := m["trace"]; ok {
		t.Error("expected 'trace' to be removed")
	}
	if m["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", m["msg"])
	}
}

func TestApply_MissingFieldNoError(t *testing.T) {
	s, _ := tagstrip.New([]string{"nonexistent"})
	input := []byte(`{"level":"warn"}`)
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "warn" {
		t.Errorf("unexpected mutation: %v", m)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s, _ := tagstrip.New([]string{"x"})
	_, err := s.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFields_ReturnsCopy(t *testing.T) {
	s, _ := tagstrip.New([]string{"a", "b"})
	fields := s.Fields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
}
