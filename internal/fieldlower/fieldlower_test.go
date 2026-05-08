package fieldlower

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
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_LowercasesField(t *testing.T) {
	l, err := New("level")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out := l.Apply(`{"level":"ERROR","msg":"oops"}`)
	m := decode(t, out)
	if got := m["level"]; got != "error" {
		t.Fatalf("expected 'error', got %q", got)
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	l, _ := New("level")
	input := `{"msg":"hello"}`
	out := l.Apply(input)
	if out != input {
		t.Fatalf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	l, _ := New("code")
	input := `{"code":404}`
	out := l.Apply(input)
	m := decode(t, out)
	if got, ok := m["code"].(float64); !ok || got != 404 {
		t.Fatalf("expected numeric 404 unchanged, got %v", m["code"])
	}
}

func TestApply_InvalidJSON_ReturnsOriginal(t *testing.T) {
	l, _ := New("level")
	input := `not-json`
	out := l.Apply(input)
	if out != input {
		t.Fatalf("expected original line returned, got %q", out)
	}
}

func TestApply_AlreadyLowercase_Unchanged(t *testing.T) {
	l, _ := New("level")
	input := `{"level":"info"}`
	out := l.Apply(input)
	m := decode(t, out)
	if got := m["level"]; got != "info" {
		t.Fatalf("expected 'info', got %q", got)
	}
}
