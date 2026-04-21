package replay

import (
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	r := New()
	if err := reg.Register("demo", r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := reg.Get("demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != r {
		t.Fatal("returned replayer does not match registered one")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := NewRegistry()
	r := New()
	_ = reg.Register("demo", r)
	if err := reg.Register("demo", r); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := NewRegistry()
	_, err := reg.Get("nope")
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Register("a", New())
	_ = reg.Register("b", New())
	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
