package tagstrip_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/tagstrip"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := tagstrip.NewRegistry()
	if err := r.Register("strip-debug", []string{"debug", "trace"}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	s, err := r.Get("strip-debug")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Stripper")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := tagstrip.NewRegistry()
	_ = r.Register("s", []string{"x"})
	if err := r.Register("s", []string{"y"}); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := tagstrip.NewRegistry()
	_, err := r.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := tagstrip.NewRegistry()
	_ = r.Register("a", []string{"f1"})
	_ = r.Register("b", []string{"f2"})
	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_InvalidFields(t *testing.T) {
	r := tagstrip.NewRegistry()
	if err := r.Register("bad", []string{}); err == nil {
		t.Fatal("expected error for empty fields slice")
	}
}
