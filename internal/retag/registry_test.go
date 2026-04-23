package retag_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/retag"
)

func makeRetagger(t *testing.T) *retag.Retagger {
	t.Helper()
	r, err := retag.New("level", map[string]string{"warn": "warning"})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := retag.NewRegistry()
	r := makeRetagger(t)
	if err := reg.Register("lvl", r); err != nil {
		t.Fatal(err)
	}
	got, err := reg.Get("lvl")
	if err != nil {
		t.Fatal(err)
	}
	if got != r {
		t.Fatal("expected same Retagger instance")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := retag.NewRegistry()
	r := makeRetagger(t)
	_ = reg.Register("lvl", r)
	if err := reg.Register("lvl", r); err == nil {
		t.Fatal("expected error on duplicate registration")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := retag.NewRegistry()
	_, err := reg.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := retag.NewRegistry()
	_ = reg.Register("a", makeRetagger(t))
	_ = reg.Register("b", makeRetagger(t))
	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
