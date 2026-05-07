package fieldround

import (
	"testing"
)

func makeRounder(t *testing.T) *Rounder {
	t.Helper()
	r, err := New("score", 2)
	if err != nil {
		t.Fatalf("makeRounder: %v", err)
	}
	return r
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register("r1", makeRounder(t)); err != nil {
		t.Fatalf("Register: %v", err)
	}
	r, err := reg.Get("r1")
	if err != nil || r == nil {
		t.Fatalf("Get: %v", err)
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Register("r1", makeRounder(t))
	if err := reg.Register("r1", makeRounder(t)); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := NewRegistry()
	_, err := reg.Get("nope")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Register("a", makeRounder(t))
	_ = reg.Register("b", makeRounder(t))
	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
