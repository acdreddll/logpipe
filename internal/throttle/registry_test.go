package throttle

import (
	"testing"
	"time"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	if err := r.Register("slow", time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	th, err := r.Get("slow")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	r.Register("dup", time.Second)
	if err := r.Register("dup", time.Second); err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing throttler")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	r.Register("a", time.Second)
	r.Register("b", time.Second)
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_InvalidCooldown(t *testing.T) {
	r := NewRegistry()
	if err := r.Register("bad", 0); err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}
