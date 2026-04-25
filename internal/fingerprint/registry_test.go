package fingerprint_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/fingerprint"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := fingerprint.NewRegistry()
	fp, _ := fingerprint.New()
	if err := reg.Register("default", fp); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, err := reg.Get("default")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != fp {
		t.Fatal("expected same Fingerprinter pointer")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := fingerprint.NewRegistry()
	fp, _ := fingerprint.New()
	_ = reg.Register("dup", fp)
	if err := reg.Register("dup", fp); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := fingerprint.NewRegistry()
	if _, err := reg.Get("missing"); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := fingerprint.NewRegistry()
	for _, n := range []string{"beta", "alpha", "gamma"} {
		fp, _ := fingerprint.New()
		_ = reg.Register(n, fp)
	}
	names := reg.Names()
	expected := []string{"alpha", "beta", "gamma"}
	if len(names) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, names)
	}
	for i, n := range names {
		if n != expected[i] {
			t.Fatalf("index %d: expected %q, got %q", i, expected[i], n)
		}
	}
}
