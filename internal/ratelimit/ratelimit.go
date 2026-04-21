// Package ratelimit provides token-bucket rate limiting for log lines.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter controls how many log lines pass through per second.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
}

// New creates a Limiter that allows up to ratePerSec log lines per second.
// A burst size equal to ratePerSec is used.
func New(ratePerSec float64) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, fmt.Errorf("ratelimit: rate must be positive, got %v", ratePerSec)
	}
	return &Limiter{
		tokens:   ratePerSec,
		max:      ratePerSec,
		rate:     ratePerSec,
		lastTick: time.Now(),
	}, nil
}

// Allow returns true if the log line should be allowed through.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}
	return false
}

// Rate returns the configured rate per second.
func (l *Limiter) Rate() float64 {
	return l.rate
}

// Tokens returns the current number of available tokens.
// This is useful for monitoring and testing purposes.
func (l *Limiter) Tokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()

	tokens := l.tokens + elapsed*l.rate
	if tokens > l.max {
		tokens = l.max
	}
	return tokens
}
