package debounce

import (
	"testing"
	"time"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	if err := r.Register("db", "msg", time.Second); err != nil {
		t.Fatal(err)
	}
	d, err := r.Get("db")
	if err != nil {
		t.Fatal(err)
	}
	if d == nil {
		t.Error("expected non-nil Debouncer")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	r.Register("db", "msg", time.Second) //nolint
	if err := r.Register("db", "msg", time.Second); err == nil {
		t.Error("expected error on duplicate registration")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("missing")
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	r.Register("a", "msg", time.Second) //nolint
	r.Register("b", "msg", time.Second) //nolint
	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_InvalidCooldown(t *testing.T) {
	r := NewRegistry()
	if err := r.Register("x", "msg", -time.Second); err == nil {
		t.Error("expected error for negative window")
	}
}
