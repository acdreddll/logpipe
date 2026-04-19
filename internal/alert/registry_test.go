package alert

import (
	"testing"
)

func makeAlert(t *testing.T, name string, fired *bool) *Alert {
	t.Helper()
	a, err := New(name, Condition{Field: "level", Operator: "eq", Value: "error"}, func(_, _ string) { *fired = true })
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return a
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	fired := false
	a := makeAlert(t, "a1", &fired)
	if err := r.Register(a); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, err := r.Get("a1")
	if err != nil || got.Name != "a1" {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	fired := false
	a := makeAlert(t, "dup", &fired)
	_ = r.Register(a)
	if err := r.Register(a); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Get("nope"); err == nil {
		t.Fatal("expected error for missing alert")
	}
}

func TestRegistry_EvaluateAll_Fires(t *testing.T) {
	r := NewRegistry()
	count := 0
	for _, name := range []string{"x", "y"} {
		n := name
		a, _ := New(n, Condition{Field: "level", Operator: "eq", Value: "error"}, func(_, _ string) { count++ })
		_ = r.Register(a)
	}
	errs := r.EvaluateAll(`{"level":"error"}`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if count != 2 {
		t.Fatalf("expected 2 fires, got %d", count)
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	fired := false
	for _, n := range []string{"p", "q"} {
		_ = r.Register(makeAlert(t, n, &fired))
	}
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
