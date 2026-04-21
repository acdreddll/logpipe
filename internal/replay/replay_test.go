package replay

import (
	"strings"
	"testing"
	"time"
)

func collect(r *Replayer, input string) []string {
	ch := make(chan string, 32)
	go func() {
		r.Run(strings.NewReader(input), ch)
		close(ch)
	}()
	var out []string
	for line := range ch {
		out = append(out, line)
	}
	return out
}

func TestRun_AllLines(t *testing.T) {
	r := New()
	lines := collect(r, "line1\nline2\nline3\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestRun_EmptyInput(t *testing.T) {
	r := New()
	lines := collect(r, "")
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestRun_MaxLines(t *testing.T) {
	r := New(WithMaxLines(2))
	lines := collect(r, "a\nb\nc\nd\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestRun_SkipsBlankLines(t *testing.T) {
	r := New()
	lines := collect(r, "a\n\nb\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestRun_WithDelay(t *testing.T) {
	r := New(WithDelay(10 * time.Millisecond))
	start := time.Now()
	collect(r, "x\ny\n")
	elapsed := time.Since(start)
	// Two lines with 10ms delay each => at least 10ms total
	if elapsed < 10*time.Millisecond {
		t.Fatalf("expected delay, elapsed=%v", elapsed)
	}
}
