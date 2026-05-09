package fieldsuffix

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", ".gz")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptySuffix(t *testing.T) {
	_, err := New("file", "")
	if err == nil {
		t.Fatal("expected error for empty suffix")
	}
}

func TestApply_AppendsSuffix(t *testing.T) {
	s, err := New("msg", "!")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := s.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["msg"]; got != "hello!" {
		t.Errorf("expected \"hello!\", got %q", got)
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	s, _ := New("msg", "!")
	const line = `{"level":"info"}`
	out, err := s.Apply(line)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	s, _ := New("code", "_suffix")
	const line = `{"code":42}`
	out, err := s.Apply(line)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got, ok := m["code"].(float64); !ok || got != 42 {
		t.Errorf("expected code=42 unchanged, got %v", m["code"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s, _ := New("msg", "!")
	_, err := s.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_MultipleFields_OnlyTargetChanged(t *testing.T) {
	s, _ := New("name", "-v2")
	out, err := s.Apply(`{"name":"logpipe","level":"debug"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["name"]; got != "logpipe-v2" {
		t.Errorf("name: expected \"logpipe-v2\", got %q", got)
	}
	if got := m["level"]; got != "debug" {
		t.Errorf("level: expected \"debug\", got %q", got)
	}
}
