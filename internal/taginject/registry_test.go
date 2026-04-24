package taginject

import (
	"testing"
)

func makeInjector(t *testing.T, value string) *Injector {
	t.Helper()
	inj, err := New(value)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return inj
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	inj := makeInjector(t, "prod")
	if err := r.Register("prod", inj); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, err := r.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != inj {
		t.Fatal("expected same injector instance")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	inj := makeInjector(t, "prod")
	_ = r.Register("prod", inj)
	if err := r.Register("prod", inj); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("beta", makeInjector(t, "beta"))
	_ = r.Register("alpha", makeInjector(t, "alpha"))
	names := r.Names()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("unexpected names: %v", names)
	}
}
