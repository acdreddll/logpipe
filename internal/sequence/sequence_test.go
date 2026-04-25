package sequence

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
	s, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.field != defaultField {
		t.Errorf("want %q, got %q", defaultField, s.field)
	}
}

func TestNew_CustomField(t *testing.T) {
	s, err := New(WithField("seq_no"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.field != "seq_no" {
		t.Errorf("want %q, got %q", "seq_no", s.field)
	}
}

func TestApply_InjectsSequenceNumber(t *testing.T) {
	s, _ := New()
	out, err := s.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m[defaultField] == nil {
		t.Fatalf("expected field %q to be present", defaultField)
	}
	if got := uint64(m[defaultField].(float64)); got != 1 {
		t.Errorf("want seq=1, got %d", got)
	}
}

func TestApply_Increments(t *testing.T) {
	s, _ := New()
	for i := 1; i <= 5; i++ {
		out, err := s.Apply([]byte(`{"x":1}`))
		if err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
		m := decode(t, out)
		got := uint64(m[defaultField].(float64))
		if got != uint64(i) {
			t.Errorf("step %d: want %d, got %d", i, i, got)
		}
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s, _ := New()
	_, err := s.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	// counter should still have advanced
	if s.Current() != 1 {
		t.Errorf("want counter=1 after failed apply, got %d", s.Current())
	}
}

func TestApply_CustomField(t *testing.T) {
	s, _ := New(WithField("n"))
	out, err := s.Apply([]byte(`{"level":"info"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["n"] == nil {
		t.Fatal("expected field 'n' to be present")
	}
}

func TestCurrent_ZeroBeforeApply(t *testing.T) {
	s, _ := New()
	if s.Current() != 0 {
		t.Errorf("want 0 before any Apply, got %d", s.Current())
	}
}
