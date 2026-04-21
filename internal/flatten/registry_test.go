package flatten

import (
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	f := New()
	if err := r.Register("default", f); err != nil {
		t.Fatal(err)
	}
	got, err := r.Get("default")
	if err != nil {
		t.Fatal(err)
	}
	if got != f {
		t.Fatal("returned wrong flattener")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("dup", New())
	if err := r.Register("dup", New()); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Get("missing"); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("b", New())
	_ = r.Register("a", New())
	names := r.Names()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("unexpected names: %v", names)
	}
}
