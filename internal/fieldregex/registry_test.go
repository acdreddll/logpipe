package fieldregex

import (
	"testing"
)

func makeReplacer(t *testing.T) *Replacer {
	t.Helper()
	r, err := New("msg", `\d+`, "NUM")
	if err != nil {
		t.Fatalf("makeReplacer: %v", err)
	}
	return r
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	r := makeReplacer(t)
	if err := reg.Register("digits", r); err != nil {
		t.Fatal(err)
	}
	got, err := reg.Get("digits")
	if err != nil {
		t.Fatal(err)
	}
	if got != r {
		t.Fatal("returned wrong replacer")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := NewRegistry()
	r := makeReplacer(t)
	_ = reg.Register("digits", r)
	if err := reg.Register("digits", r); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := NewRegistry()
	if _, err := reg.Get("nope"); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Register("a", makeReplacer(t))
	_ = reg.Register("b", makeReplacer(t))
	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
