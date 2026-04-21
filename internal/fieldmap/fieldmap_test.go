package fieldmap

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyFrom(t *testing.T) {
	_, err := New([]Mapping{{From: "", To: "dst"}})
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestNew_EmptyTo(t *testing.T) {
	_, err := New([]Mapping{{From: "src", To: ""}})
	if err == nil {
		t.Fatal("expected error for empty To")
	}
}

func TestApply_CopyField(t *testing.T) {
	m, _ := New([]Mapping{{From: "msg", To: "message", Delete: false}})
	out, err := m.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if obj["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", obj["message"])
	}
	if _, ok := obj["msg"]; !ok {
		t.Error("expected original msg field to remain")
	}
}

func TestApply_MoveField(t *testing.T) {
	m, _ := New([]Mapping{{From: "msg", To: "message", Delete: true}})
	out, err := m.Apply([]byte(`{"msg":"hello","level":"info"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if obj["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", obj["message"])
	}
	if _, ok := obj["msg"]; ok {
		t.Error("expected original msg field to be deleted")
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	m, _ := New([]Mapping{{From: "nonexistent", To: "dst", Delete: true}})
	input := []byte(`{"level":"warn"}`)
	out, err := m.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if _, ok := obj["dst"]; ok {
		t.Error("dst should not be present when source is missing")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m, _ := New([]Mapping{{From: "a", To: "b"}})
	_, err := m.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
