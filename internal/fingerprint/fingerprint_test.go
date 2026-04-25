package fingerprint_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/fingerprint"
)

func mustNew(t *testing.T, opts ...fingerprint.Option) *fingerprint.Fingerprinter {
	t.Helper()
	fp, err := fingerprint.New(opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return fp
}

func TestCompute_WholeEvent_Deterministic(t *testing.T) {
	fp := mustNew(t)
	line := []byte(`{"level":"info","msg":"hello"}`)
	a, err := fp.Compute(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := fp.Compute(line)
	if a != b {
		t.Fatalf("expected same fingerprint, got %q vs %q", a, b)
	}
}

func TestCompute_DifferentEvents_DifferentFingerprints(t *testing.T) {
	fp := mustNew(t)
	a, _ := fp.Compute([]byte(`{"msg":"hello"}`))
	b, _ := fp.Compute([]byte(`{"msg":"world"}`))
	if a == b {
		t.Fatal("expected different fingerprints for different events")
	}
}

func TestCompute_WithFields_IgnoresOtherFields(t *testing.T) {
	fp := mustNew(t, fingerprint.WithFields("level"))
	a, _ := fp.Compute([]byte(`{"level":"error","msg":"a"}`))
	b, _ := fp.Compute([]byte(`{"level":"error","msg":"b"}`))
	if a != b {
		t.Fatalf("expected same fingerprint when restricted field is equal, got %q vs %q", a, b)
	}
}

func TestCompute_WithFields_DifferentValues(t *testing.T) {
	fp := mustNew(t, fingerprint.WithFields("level"))
	a, _ := fp.Compute([]byte(`{"level":"info"}`))
	b, _ := fp.Compute([]byte(`{"level":"error"}`))
	if a == b {
		t.Fatal("expected different fingerprints for different field values")
	}
}

func TestCompute_KeyOrderIndependent(t *testing.T) {
	fp := mustNew(t)
	a, _ := fp.Compute([]byte(`{"a":"1","b":"2"}`))
	b, _ := fp.Compute([]byte(`{"b":"2","a":"1"}`))
	if a != b {
		t.Fatalf("expected key-order-independent fingerprint, got %q vs %q", a, b)
	}
}

func TestCompute_InvalidJSON_ReturnsError(t *testing.T) {
	fp := mustNew(t)
	_, err := fp.Compute([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
