package sampling_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/sampling"
)

func TestNew_RateOne_KeepsAll(t *testing.T) {
	s := sampling.New(1)
	for i := 0; i < 100; i++ {
		if !s.Sample() {
			t.Fatal("rate=1 sampler dropped a line")
		}
	}
}

func TestNew_RateZero_TreatedAsOne(t *testing.T) {
	s := sampling.New(0)
	for i := 0; i < 10; i++ {
		if !s.Sample() {
			t.Fatal("rate=0 sampler dropped a line")
		}
	}
}

func TestNew_RateN_KeepsEveryNth(t *testing.T) {
	s := sampling.New(4)
	kept := 0
	for i := 0; i < 100; i++ {
		if s.Sample() {
			kept++
		}
	}
	if kept != 25 {
		t.Fatalf("expected 25 kept, got %d", kept)
	}
}

func TestNewRandom_ProbOne_KeepsAll(t *testing.T) {
	s := sampling.NewRandom(1.0)
	for i := 0; i < 50; i++ {
		if !s.Sample() {
			t.Fatal("prob=1.0 sampler dropped a line")
		}
	}
}

func TestNewRandom_ProbZero_DropsAll(t *testing.T) {
	s := sampling.NewRandom(0.0)
	for i := 0; i < 50; i++ {
		if s.Sample() {
			t.Fatal("prob=0.0 sampler kept a line")
		}
	}
}

func TestNewRandom_Prob_ApproxRate(t *testing.T) {
	s := sampling.NewRandom(0.5)
	kept := 0
	total := 10000
	for i := 0; i < total; i++ {
		if s.Sample() {
			kept++
		}
	}
	ratio := float64(kept) / float64(total)
	if ratio < 0.45 || ratio > 0.55 {
		t.Fatalf("expected ~0.5 keep ratio, got %.3f", ratio)
	}
}
