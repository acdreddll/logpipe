package fieldtrim

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
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

func TestApply_TrimsWhitespace(t *testing.T) {
	tr, err := New("msg")
	if err != nil {
		t.Fatal(err)
	}
	out, err := tr.Apply(`{"msg":"  hello world  "}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "hello world" {
		t.Errorf("got %q, want %q", m["msg"], "hello world")
	}
}

func TestApply_CustomCutset(t *testing.T) {
	tr, err := New("msg", WithCutset("-"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := tr.Apply(`{"msg":"---hello---"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("got %q, want %q", m["msg"], "hello")
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	tr, _ := New("msg")
	const input = `{"level":"info"}`
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if out != input {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	tr, _ := New("count")
	const input = `{"count":42}`
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["count"].(float64) != 42 {
		t.Errorf("expected count=42, got %v", m["count"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr, _ := New("msg")
	_, err := tr.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_AlreadyTrimmed_Unchanged(t *testing.T) {
	tr, _ := New("msg")
	const input = `{"msg":"clean"}`
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "clean" {
		t.Errorf("got %q, want %q", m["msg"], "clean")
	}
}
