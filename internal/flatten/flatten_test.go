package flatten

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

func TestApply_FlatInput(t *testing.T) {
	f := New()
	out, err := f.Apply(`{"level":"info","msg":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["level"] != "info" || m["msg"] != "ok" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestApply_NestedObject(t *testing.T) {
	f := New()
	out, err := f.Apply(`{"a":{"b":{"c":42}}}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["a.b.c"] != float64(42) {
		t.Fatalf("expected a.b.c=42, got %v", m)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := New(WithSeparator("_"))
	out, err := f.Apply(`{"x":{"y":"z"}}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["x_y"] != "z" {
		t.Fatalf("expected x_y=z, got %v", m)
	}
}

func TestApply_WithPrefix(t *testing.T) {
	f := New(WithPrefix("log"))
	out, err := f.Apply(`{"level":"warn"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["log.level"] != "warn" {
		t.Fatalf("expected log.level=warn, got %v", m)
	}
}

func TestApply_ArrayPreserved(t *testing.T) {
	f := New()
	out, err := f.Apply(`{"tags":["a","b"]}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["tags"]; !ok {
		t.Fatal("expected tags key to be preserved")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	f := New()
	_, err := f.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
