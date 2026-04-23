package timestamp

import (
	"encoding/json"
	"testing"
	"time"
)

func decode(t *testing.T, line []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
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

func TestApply_InjectsFieldWhenAbsent(t *testing.T) {
	p, _ := New("ts")
	out, err := p.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["ts"]; !ok {
		t.Error("expected ts field to be injected")
	}
}

func TestApply_InjectNow_OverwritesExisting(t *testing.T) {
	p, _ := New("ts", WithInjectNow())
	before := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	input, _ := json.Marshal(map[string]interface{}{"ts": before})
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["ts"] == before {
		t.Error("expected ts to be overwritten with current time")
	}
}

func TestApply_ReformatsExistingTimestamp(t *testing.T) {
	p, _ := New("ts", WithFormat(time.RFC3339))
	input := []byte(`{"ts":"2024-01-15T10:00:00Z"}`)
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["ts"] != "2024-01-15T10:00:00Z" {
		t.Errorf("unexpected ts value: %v", m["ts"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := New("ts")
	_, err := p.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_UnparsableTimestamp_LeftUntouched(t *testing.T) {
	p, _ := New("ts", WithFormat(time.RFC3339))
	input := []byte(`{"ts":"not-a-date"}`)
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["ts"] != "not-a-date" {
		t.Errorf("expected original value preserved, got %v", m["ts"])
	}
}
