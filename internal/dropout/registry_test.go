package dropout

import (
	"testing"
)

func makeDropper(t *testing.T) *Dropper {
	t.Helper()
	d, err := New("level", []string{"debug"})
	if err != nil {
		t.Fatalf("makeDropper: %v", err)
	}
	return d
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	if err := r.Register("noisy", makeDropper(t)); err != nil {
		t.Fatalf("Register: %v", err)
	}
	d, err := r.Get("noisy")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil dropper")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("noisy", makeDropper(t))
	if err := r.Register("noisy", makeDropper(t)); err == nil {
		t.Fatal("expected error on duplicate registration")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("ghost")
	if err == nil {
		t.Fatal("expected error for missing dropper")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("a", makeDropper(t))
	_ = r.Register("b", makeDropper(t))
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
