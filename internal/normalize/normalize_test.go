package normalize

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

func TestNew_NoOptions(t *testing.T) {
	_, err := New()
	if err == nil {
		t.Fatal("expected error with no options")
	}
}

func TestApply_Lowercase(t *testing.T) {
	n, err := New(WithLowercase())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := n.Apply(`{"Level":"info","MSG":"hello"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["level"]; !ok {
		t.Error("expected key 'level'")
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected key 'msg'")
	}
	if _, ok := m["Level"]; ok {
		t.Error("original key 'Level' should be gone")
	}
}

func TestApply_ExplicitMapping(t *testing.T) {
	n, err := New(WithMapping(map[string]string{"log_level": "level", "message": "msg"}))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := n.Apply(`{"log_level":"warn","message":"oops","ts":"2024-01-01"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", m["level"])
	}
	if m["msg"] != "oops" {
		t.Errorf("expected msg=oops, got %v", m["msg"])
	}
	if _, ok := m["ts"]; !ok {
		t.Error("unmapped key 'ts' should be preserved")
	}
}

func TestApply_MappingTakesPrecedenceOverLowercase(t *testing.T) {
	n, err := New(WithLowercase(), WithMapping(map[string]string{"Timestamp": "time"}))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := n.Apply(`{"Timestamp":"2024-01-01","Level":"debug"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["time"]; !ok {
		t.Error("expected explicit rename 'time'")
	}
	if _, ok := m["level"]; !ok {
		t.Error("expected lowercased key 'level'")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	n, err := New(WithLowercase())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = n.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
