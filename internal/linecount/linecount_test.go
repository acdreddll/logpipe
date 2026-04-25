package linecount

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_DefaultField(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.field != defaultField {
		t.Fatalf("expected field %q, got %q", defaultField, c.field)
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New(WithField(""))
	if err == nil {
		t.Fatal("expected error for empty field, got nil")
	}
}

func TestApply_InjectsLineNumber(t *testing.T) {
	c, _ := New()
	out, err := c.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if v, ok := m[defaultField]; !ok {
		t.Fatalf("field %q missing from output", defaultField)
	} else if int64(v.(float64)) != 1 {
		t.Fatalf("expected line 1, got %v", v)
	}
}

func TestApply_Increments(t *testing.T) {
	c, _ := New()
	for i := int64(1); i <= 5; i++ {
		out, err := c.Apply([]byte(`{"x":1}`))
		if err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
		m := decode(t, out)
		got := int64(m[defaultField].(float64))
		if got != i {
			t.Fatalf("step %d: expected %d, got %d", i, i, got)
		}
	}
}

func TestApply_CustomField(t *testing.T) {
	c, _ := New(WithField("seq"))
	out, _ := c.Apply([]byte(`{"msg":"test"}`))
	m := decode(t, out)
	if _, ok := m["seq"]; !ok {
		t.Fatal("custom field 'seq' not found in output")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := New()
	_, err := c.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	c, _ := New()
	c.Apply([]byte(`{"a":1}`))
	c.Apply([]byte(`{"a":2}`))
	if c.Value() != 2 {
		t.Fatalf("expected value 2 before reset, got %d", c.Value())
	}
	c.Reset()
	if c.Value() != 0 {
		t.Fatalf("expected value 0 after reset, got %d", c.Value())
	}
	out, _ := c.Apply([]byte(`{"a":3}`))
	m := decode(t, out)
	if int64(m[defaultField].(float64)) != 1 {
		t.Fatal("expected counter to restart at 1 after reset")
	}
}
