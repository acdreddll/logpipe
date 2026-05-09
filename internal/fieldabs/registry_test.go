package fieldabs_test

import (
	"testing"

	"logpipe/internal/fieldabs"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := fieldabs.NewRegistry()
	if err := r.Register("abs_latency", "latency"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	p, err := r.Get("abs_latency")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil processor")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := fieldabs.NewRegistry()
	_ = r.Register("p", "value")
	if err := r.Register("p", "value"); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := fieldabs.NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing processor")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := fieldabs.NewRegistry()
	_ = r.Register("a", "x")
	_ = r.Register("b", "y")
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_InvalidField(t *testing.T) {
	r := fieldabs.NewRegistry()
	if err := r.Register("bad", ""); err == nil {
		t.Fatal("expected error for empty field name")
	}
}
