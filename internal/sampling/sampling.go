// Package sampling provides log sampling strategies for rate-limiting
// high-volume log streams before they reach the pipeline.
package sampling

import (
	"math/rand"
	"sync/atomic"
)

// Sampler decides whether a given log line should be kept.
type Sampler interface {
	Sample() bool
}

// RateSampler keeps approximately 1/N log lines.
type RateSampler struct {
	rate    int
	counter uint64
}

// New returns a Sampler that retains roughly 1 out of every rate lines.
// A rate of 1 keeps all lines. A rate <= 0 is treated as 1.
func New(rate int) Sampler {
	if rate <= 0 {
		rate = 1
	}
	if rate == 1 {
		return &passThroughSampler{}
	}
	return &RateSampler{rate: rate}
}

// Sample returns true if the line should be kept.
func (s *RateSampler) Sample() bool {
	n := atomic.AddUint64(&s.counter, 1)
	return n%uint64(s.rate) == 0
}

// RandomSampler keeps each line with probability p in [0.0, 1.0].
type RandomSampler struct {
	prob float64
}

// NewRandom returns a Sampler that keeps each line with probability p.
func NewRandom(p float64) Sampler {
	if p >= 1.0 {
		return &passThroughSampler{}
	}
	if p <= 0.0 {
		return &dropAllSampler{}
	}
	return &RandomSampler{prob: p}
}

func (s *RandomSampler) Sample() bool {
	return rand.Float64() < s.prob
}

type passThroughSampler struct{}

func (p *passThroughSampler) Sample() bool { return true }

type dropAllSampler struct{}

func (d *dropAllSampler) Sample() bool { return false }
