package coalesce

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

func TestNew_EmptyDest(t *testing.T) {
	_, err := New("", []string{"a"})
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestNew_NoSources(t *testing.T) {
	_, err := New("out", nil)
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestApply_FirstSourceWins(t *testing.T) {
	c, _ := New("result", []string{"a", "b", "c"})
	out, err := c.Apply(`{"a":"hello","b":"world"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["result"] != "hello" {
		t.Errorf("expected 'hello', got %v", m["result"])
	}
}

func TestApply_SkipsEmptyString(t *testing.T) {
	c, _ := New("result", []string{"a", "b"})
	out, err := c.Apply(`{"a":"","b":"fallback"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["result"] != "fallback" {
		t.Errorf("expected 'fallback', got %v", m["result"])
	}
}

func TestApply_SkipsMissingField(t *testing.T) {
	c, _ := New("result", []string{"missing", "b"})
	out, err := c.Apply(`{"b":"found"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["result"] != "found" {
		t.Errorf("expected 'found', got %v", m["result"])
	}
}

func TestApply_NoMatchReturnsUnchanged(t *testing.T) {
	c, _ := New("result", []string{"x", "y"})
	input := `{"a":"value"}`
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected unchanged line, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := New("result", []string{"a"})
	_, err := c.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
