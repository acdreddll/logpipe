package aggregator

import (
	"testing"
)

func mustNew(t *testing.T, field string, op Op) *Aggregator {
	t.Helper()
	a, err := New(field, op)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return a
}

func TestNew_InvalidOp(t *testing.T) {
	_, err := New("field", Op("median"))
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", OpCount)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestAdd_Count(t *testing.T) {
	a := mustNew(t, "level", OpCount)
	for i := 0; i < 5; i++ {
		if err := a.Add([]byte(`{"level":"info"}`)); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}
	if got := a.Result(); got != 5 {
		t.Errorf("expected 5, got %v", got)
	}
}

func TestAdd_Sum(t *testing.T) {
	a := mustNew(t, "duration", OpSum)
	lines := []string{
		`{"duration":10}`,
		`{"duration":20}`,
		`{"duration":30}`,
	}
	for _, l := range lines {
		a.Add([]byte(l))
	}
	if got := a.Result(); got != 60 {
		t.Errorf("expected 60, got %v", got)
	}
}

func TestAdd_Min(t *testing.T) {
	a := mustNew(t, "latency", OpMin)
	for _, v := range []string{`{"latency":5}`, `{"latency":2}`, `{"latency":9}`} {
		a.Add([]byte(v))
	}
	if got := a.Result(); got != 2 {
		t.Errorf("expected 2, got %v", got)
	}
}

func TestAdd_Max(t *testing.T) {
	a := mustNew(t, "latency", OpMax)
	for _, v := range []string{`{"latency":5}`, `{"latency":2}`, `{"latency":9}`} {
		a.Add([]byte(v))
	}
	if got := a.Result(); got != 9 {
		t.Errorf("expected 9, got %v", got)
	}
}

func TestAdd_InvalidJSON(t *testing.T) {
	a := mustNew(t, "x", OpSum)
	if err := a.Add([]byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestAdd_MissingField_NoError(t *testing.T) {
	a := mustNew(t, "missing", OpSum)
	if err := a.Add([]byte(`{"other":1}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := a.Result(); got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestReset(t *testing.T) {
	a := mustNew(t, "x", OpCount)
	a.Add([]byte(`{"x":1}`))
	a.Add([]byte(`{"x":2}`))
	a.Reset()
	if got := a.Result(); got != 0 {
		t.Errorf("expected 0 after reset, got %v", got)
	}
}
