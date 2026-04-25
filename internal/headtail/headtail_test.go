package headtail

import (
	"testing"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := New(Head, 0)
	if err == nil {
		t.Fatal("expected error for limit 0")
	}
}

func TestNew_ValidLimit(t *testing.T) {
	p, err := New(Head, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil processor")
	}
}

func TestHead_KeepsFirstN(t *testing.T) {
	p, _ := New(Head, 3)
	for _, line := range []string{"a", "b", "c", "d", "e"} {
		p.Add(line)
	}
	out := p.Flush()
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
	if out[0] != "a" || out[2] != "c" {
		t.Errorf("unexpected head lines: %v", out)
	}
}

func TestTail_KeepsLastN(t *testing.T) {
	p, _ := New(Tail, 3)
	for _, line := range []string{"a", "b", "c", "d", "e"} {
		p.Add(line)
	}
	out := p.Flush()
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
	if out[0] != "c" || out[2] != "e" {
		t.Errorf("unexpected tail lines: %v", out)
	}
}

func TestFlush_ResetsBuffer(t *testing.T) {
	p, _ := New(Head, 5)
	p.Add("x")
	p.Add("y")
	p.Flush()
	if p.Len() != 0 {
		t.Fatalf("expected empty buffer after flush, got %d", p.Len())
	}
}

func TestAdd_SkipsEmptyLines(t *testing.T) {
	p, _ := New(Head, 5)
	p.Add("")
	p.Add("valid")
	if p.Len() != 1 {
		t.Fatalf("expected 1 line, got %d", p.Len())
	}
}

func TestHead_FewerThanLimit(t *testing.T) {
	p, _ := New(Head, 10)
	p.Add("only")
	out := p.Flush()
	if len(out) != 1 || out[0] != "only" {
		t.Errorf("unexpected output: %v", out)
	}
}
