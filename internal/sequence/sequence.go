// Package sequence assigns a monotonically increasing sequence number to each
// log line that passes through it. The counter is per-instance and resets when
// the Sequencer is recreated.
package sequence

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

const defaultField = "_seq"

// Sequencer injects a sequence number into every JSON log line.
type Sequencer struct {
	field   string
	counter atomic.Uint64
}

// Option is a functional option for Sequencer.
type Option func(*Sequencer)

// WithField overrides the JSON field name used for the sequence number.
func WithField(field string) Option {
	return func(s *Sequencer) {
		if field != "" {
			s.field = field
		}
	}
}

// New creates a Sequencer that starts counting from 1.
func New(opts ...Option) (*Sequencer, error) {
	s := &Sequencer{field: defaultField}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Apply injects the next sequence number into the given JSON line.
// The counter increments even when the line cannot be parsed so that sequence
// numbers remain unique across the stream.
func (s *Sequencer) Apply(line []byte) ([]byte, error) {
	n := s.counter.Add(1)

	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return line, fmt.Errorf("sequence: invalid JSON: %w", err)
	}

	obj[s.field] = n

	out, err := json.Marshal(obj)
	if err != nil {
		return line, fmt.Errorf("sequence: marshal: %w", err)
	}
	return out, nil
}

// Current returns the last sequence number that was assigned.
// It returns 0 if Apply has never been called.
func (s *Sequencer) Current() uint64 {
	return s.counter.Load()
}
