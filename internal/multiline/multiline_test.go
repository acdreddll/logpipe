package multiline

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]string {
	t.Helper()
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "message")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[", "message")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_DefaultField(t *testing.T) {
	a, err := New(`^\d{4}`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.field != "message" {
		t.Errorf("expected field=message, got %q", a.field)
	}
}

func TestAdd_SingleEvent(t *testing.T) {
	a, _ := New(`^\d{4}`, "message")

	out, flushed, err := a.Add("2024-01-01 first line")
	if err != nil || flushed || out != nil {
		t.Fatal("first line should buffer without flushing")
	}

	out, flushed, err = a.Add("  continuation")
	if err != nil || flushed || out != nil {
		t.Fatal("continuation should buffer without flushing")
	}

	out, err = a.Flush()
	if err != nil {
		t.Fatalf("flush error: %v", err)
	}
	m := decode(t, out)
	if m["message"] != "2024-01-01 first line\n  continuation" {
		t.Errorf("unexpected message: %q", m["message"])
	}
}

func TestAdd_TwoEvents(t *testing.T) {
	a, _ := New(`^\d{4}`, "message")

	a.Add("2024-01-01 event one") //nolint
	a.Add("  stack line")         //nolint

	out, flushed, err := a.Add("2024-01-02 event two")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !flushed {
		t.Fatal("expected flush on second start line")
	}
	m := decode(t, out)
	if m["message"] != "2024-01-01 event one\n  stack line" {
		t.Errorf("first event wrong: %q", m["message"])
	}

	out, _ = a.Flush()
	m = decode(t, out)
	if m["message"] != "2024-01-02 event two" {
		t.Errorf("second event wrong: %q", m["message"])
	}
}

func TestFlush_EmptyBuffer(t *testing.T) {
	a, _ := New(`^\d{4}`, "message")
	out, err := a.Flush()
	if err != nil || out != nil {
		t.Fatal("flush on empty buffer should return nil, nil")
	}
}
