package eventsource

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("app", WithField(""))
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := New("nginx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name() != "nginx" {
		t.Errorf("Name() = %q, want %q", s.Name(), "nginx")
	}
	if s.Field() != defaultField {
		t.Errorf("Field() = %q, want %q", s.Field(), defaultField)
	}
}

func TestApply_InjectsSource(t *testing.T) {
	s, _ := New("myapp")
	out, err := s.Apply([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["source"]; got != "myapp" {
		t.Errorf("source = %v, want %q", got, "myapp")
	}
}

func TestApply_OverwritesExisting(t *testing.T) {
	s, _ := New("new")
	out, err := s.Apply([]byte(`{"source":"old"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["source"] != "new" {
		t.Errorf("source = %v, want %q", m["source"], "new")
	}
}

func TestApply_CustomField(t *testing.T) {
	s, _ := New("svc", WithField("origin"))
	out, err := s.Apply([]byte(`{"msg":"ok"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["origin"] != "svc" {
		t.Errorf("origin = %v, want %q", m["origin"], "svc")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s, _ := New("x")
	_, err := s.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
