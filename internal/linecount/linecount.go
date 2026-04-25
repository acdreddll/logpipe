// Package linecount provides a processor that injects a running line-count
// field into each JSON log event passing through the pipeline.
package linecount

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

const defaultField = "_line"

// Counter injects a monotonically increasing line number into each log event.
type Counter struct {
	field string
	n     atomic.Int64
}

// Option is a functional option for Counter.
type Option func(*Counter)

// WithField overrides the destination field name (default: "_line").
func WithField(field string) Option {
	return func(c *Counter) {
		c.field = field
	}
}

// New creates a new Counter. The counter starts at 1 for the first event.
func New(opts ...Option) (*Counter, error) {
	c := &Counter{field: defaultField}
	for _, o := range opts {
		o(c)
	}
	if c.field == "" {
		return nil, fmt.Errorf("linecount: field name must not be empty")
	}
	return c, nil
}

// Apply injects the current line number into the JSON object and returns the
// updated line. The counter is incremented atomically before injection so the
// first call produces field value 1.
func (c *Counter) Apply(line []byte) ([]byte, error) {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("linecount: invalid JSON: %w", err)
	}

	num := c.n.Add(1)
	obj[c.field] = num

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("linecount: marshal: %w", err)
	}
	return out, nil
}

// Reset resets the counter back to zero.
func (c *Counter) Reset() {
	c.n.Store(0)
}

// Value returns the current counter value (i.e. the number of events seen).
func (c *Counter) Value() int64 {
	return c.n.Load()
}
