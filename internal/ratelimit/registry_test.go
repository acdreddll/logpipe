package ratelimit_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/ratelimit"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := ratelimit.NewRegistry()
	if err := r.Register("route-a", 50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l, err := r.Get("route-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Rate() != 50 {
		t.Fatalf("expected rate 50, got %v", l.Rate())
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := ratelimit.NewRegistry()
	r.Register("route-b", 10)
	if err := r.Register("route-b", 20); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := ratelimit.NewRegistry()
	_, err := r.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing limiter")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := ratelimit.NewRegistry()
	r.Register("a", 1)
	r.Register("b", 2)
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_InvalidRate(t *testing.T) {
	r := ratelimit.NewRegistry()
	if err := r.Register("bad", -1); err == nil {
		t.Fatal("expected error for invalid rate")
	}
}
