package ratelimit_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/ratelimit"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(0)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
	_, err = ratelimit.New(-5)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := ratelimit.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Rate() != 10 {
		t.Fatalf("expected rate 10, got %v", l.Rate())
	}
}

func TestAllow_BurstAllowed(t *testing.T) {
	l, _ := ratelimit.New(5)
	allowed := 0
	for i := 0; i < 5; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed in burst, got %d", allowed)
	}
}

func TestAllow_ExceedBurst(t *testing.T) {
	l, _ := ratelimit.New(3)
	// drain burst
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected token to be denied after burst exhausted")
	}
}

func TestAllow_RefillOverTime(t *testing.T) {
	l, _ := ratelimit.New(100)
	// drain all
	for i := 0; i < 100; i++ {
		l.Allow()
	}
	time.Sleep(50 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected token after refill delay")
	}
}
