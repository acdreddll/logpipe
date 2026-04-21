package throttle

import (
	"testing"
	"time"
)

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative cooldown")
	}
}

func TestAllow_FirstOccurrence(t *testing.T) {
	th, _ := New(time.Second)
	line := []byte(`{"level":"error","msg":"oops"}`)
	ok, err := th.Allow(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestAllow_SuppressedWithinCooldown(t *testing.T) {
	th, _ := New(time.Hour)
	line := []byte(`{"level":"error","msg":"oops"}`)
	th.Allow(line) // prime
	ok, err := th.Allow(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected second occurrence to be suppressed")
	}
}

func TestAllow_AllowedAfterCooldown(t *testing.T) {
	th, _ := New(50 * time.Millisecond)
	base := time.Now()
	th.now = func() time.Time { return base }

	line := []byte(`{"level":"warn","msg":"slow"}`)
	th.Allow(line)

	th.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	ok, err := th.Allow(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected line to be allowed after cooldown elapsed")
	}
}

func TestAllow_InvalidJSON(t *testing.T) {
	th, _ := New(time.Second)
	_, err := th.Allow([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	th, _ := New(50 * time.Millisecond)
	base := time.Now()
	th.now = func() time.Time { return base }

	th.Allow([]byte(`{"msg":"a"}`))
	th.Allow([]byte(`{"msg":"b"}`))

	th.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	th.Purge()

	if len(th.seen) != 0 {
		t.Fatalf("expected seen map to be empty after purge, got %d entries", len(th.seen))
	}
}
