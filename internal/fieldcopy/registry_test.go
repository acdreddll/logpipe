package fieldcopy_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/fieldcopy"
)

func makeCopier(t *testing.T) *fieldcopy.Copier {
	t.Helper()
	c, err := fieldcopy.New("src", "dst")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := fieldcopy.NewRegistry()
	if err := r.Register("copy1", makeCopier(t)); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if _, err := r.Get("copy1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := fieldcopy.NewRegistry()
	_ = r.Register("copy1", makeCopier(t))
	if err := r.Register("copy1", makeCopier(t)); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := fieldcopy.NewRegistry()
	if _, err := r.Get("nope"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := fieldcopy.NewRegistry()
	_ = r.Register("b", makeCopier(t))
	_ = r.Register("a", makeCopier(t))
	names := r.Names()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected names: %v", names)
	}
}
