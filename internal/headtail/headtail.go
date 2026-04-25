// Package headtail provides a processor that keeps only the first N
// (head) or last N (tail) lines from a log stream window.
package headtail

import (
	"errors"
	"sync"
)

// Mode selects whether Head or Tail behaviour is used.
type Mode int

const (
	Head Mode = iota
	Tail
)

// Processor buffers log lines and returns only the head or tail slice.
type Processor struct {
	mu    sync.Mutex
	mode  Mode
	limit int
	buf   []string
}

// New creates a Processor with the given mode and line limit.
// limit must be >= 1.
func New(mode Mode, limit int) (*Processor, error) {
	if limit < 1 {
		return nil, errors.New("headtail: limit must be >= 1")
	}
	return &Processor{mode: mode, limit: limit}, nil
}

// Add appends a raw log line to the internal buffer.
// For Head mode the buffer is capped at limit; for Tail the oldest
// entry is evicted once capacity is exceeded.
func (p *Processor) Add(line string) {
	if line == "" {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.mode {
	case Head:
		if len(p.buf) < p.limit {
			p.buf = append(p.buf, line)
		}
	case Tail:
		p.buf = append(p.buf, line)
		if len(p.buf) > p.limit {
			p.buf = p.buf[len(p.buf)-p.limit:]
		}
	}
}

// Flush returns the accumulated lines and resets the buffer.
func (p *Processor) Flush() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]string, len(p.buf))
	copy(out, p.buf)
	p.buf = p.buf[:0]
	return out
}

// Len returns the current number of buffered lines.
func (p *Processor) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.buf)
}
