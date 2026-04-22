package rename

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyMapping(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty mapping")
	}
}

func TestNew_EmptySourceKey(t *testing.T) {
	_, err := New(map[string]string{"": "dst"})
	if err == nil {
		t.Fatal("expected error for empty source key")
	}
}

func TestNew_EmptyDestKey(t *testing.T) {
	_, err := New(map[string]string{"src": ""})
	if err == nil {
		t.Fatal("expected error for empty destination key")
	}
}

func TestApply_RenamesField(t *testing.T) {
	r, err := New(map[string]string{"msg": "message"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := r.Apply([]byte(`{"msg":"hello","level":"info"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["msg"]; ok {
		t.Error("old field 'msg' should have been removed")
	}
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestApply_MissingSourceField_NoError(t *testing.T) {
	r, _ := New(map[string]string{"missing": "present"})
	out, err := r.Apply([]byte(`{"level":"warn"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "warn" {
		t.Errorf("unexpected level: %v", m["level"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r, _ := New(map[string]string{"a": "b"})
	_, err := r.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_MultipleRenames(t *testing.T) {
	r, _ := New(map[string]string{"ts": "timestamp", "msg": "message"})
	out, err := r.Apply([]byte(`{"ts":"2024-01-01","msg":"hi","level":"debug"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["timestamp"] != "2024-01-01" {
		t.Errorf("expected timestamp=2024-01-01, got %v", m["timestamp"])
	}
	if m["message"] != "hi" {
		t.Errorf("expected message=hi, got %v", m["message"])
	}
	if _, ok := m["ts"]; ok {
		t.Error("old field 'ts' should be gone")
	}
}
