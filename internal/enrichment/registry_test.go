package enrichment

import (
	"testing"
)

func mustNew(t *testing.T, fields map[string]string) *Enricher {
	t.Helper()
	e, err := New(fields)
	if err != nil {
		t.Fatalf("mustNew: %v", err)
	}
	return e
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	e := mustNew(t, map[string]string{"env": "test"})
	if err := r.Register("default", e); err != nil {
		t.Fatalf("register: %v", err)
	}
	got, err := r.Get("default")
	if err != nil || got != e {
		t.Fatalf("get: %v %v", got, err)
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	e := mustNew(t, map[string]string{"env": "test"})
	r.Register("dup", e)
	if err := r.Register("dup", e); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Get("missing"); err == nil {
		t.Fatal("expected error for missing enricher")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	r.Register("a", mustNew(t, map[string]string{"x": "1"}))
	r.Register("b", mustNew(t, map[string]string{"y": "2"}))
	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
