package buffer

import (
	"sync"
	"testing"
	"time"
)

func collect(t *testing.T, size int, interval time.Duration, lines []string) []string {
	t.Helper()
	var mu sync.Mutex
	var got []string
	b := New(size, interval, func(batch []string) {
		mu.Lock()
		got = append(got, batch...)
		mu.Unlock()
	})
	for _, l := range lines {
		b.Add(l)
	}
	b.Stop()
	return got
}

func TestAdd_FlushOnCapacity(t *testing.T) {
	got := collect(t, 3, time.Hour, []string{"a", "b", "c", "d"})
	if len(got) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(got))
	}
}

func TestStop_FlushesRemainder(t *testing.T) {
	got := collect(t, 10, time.Hour, []string{"x", "y"})
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestAdd_IntervalFlush(t *testing.T) {
	var mu sync.Mutex
	var got []string
	b := New(100, 20*time.Millisecond, func(batch []string) {
		mu.Lock()
		got = append(got, batch...)
		mu.Unlock()
	})
	b.Add("tick")
	time.Sleep(60 * time.Millisecond)
	b.Stop()
	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected at least one line flushed by ticker")
	}
}

func TestNew_DefaultSize(t *testing.T) {
	b := New(0, time.Hour, func([]string) {})
	if b.cap != 100 {
		t.Fatalf("expected default cap 100, got %d", b.cap)
	}
	b.Stop()
}

func TestStop_EmptyBuffer(t *testing.T) {
	flushed := false
	b := New(10, time.Hour, func(batch []string) { flushed = true })
	b.Stop()
	if flushed {
		t.Fatal("expected no flush for empty buffer")
	}
}
