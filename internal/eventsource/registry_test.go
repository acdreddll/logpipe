package eventsource

import (
	"testing"
)

func makeSource(t *testing.T, name string) *Source {
	t.Helper()
	s, err := New(name)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	s := makeSource(t, "app")
	if err := r.Register("app", s); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, err := r.Get("app")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name() != "app" {
		t.Errorf("Name() = %q, want %q", got.Name(), "app")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	s := makeSource(t, "app")
	_ = r.Register("app", s)
	if err := r.Register("app", s); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Get("missing"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	for _, n := range []string{"b", "a", "c"} {
		_ = r.Register(n, makeSource(t, n))
	}
	names := r.Names()
	want := []string{"a", "b", "c"}
	for i, w := range want {
		if names[i] != w {
			t.Errorf("Names()[%d] = %q, want %q", i, names[i], w)
		}
	}
}
