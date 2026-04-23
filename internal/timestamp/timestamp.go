// Package timestamp provides a processor that normalises or injects a
// timestamp field into structured log lines.
package timestamp

import (
	"encoding/json"
	"fmt"
	"time"
)

// Processor rewrites or injects a timestamp field in a JSON log line.
type Processor struct {
	field  string
	format string
	injectNow bool
}

// Option is a functional option for Processor.
type Option func(*Processor)

// WithFormat sets the time format used when parsing an existing field.
// Defaults to time.RFC3339.
func WithFormat(format string) Option {
	return func(p *Processor) {
		p.format = format
	}
}

// WithInjectNow causes the processor to overwrite the field with the
// current UTC time regardless of any existing value.
func WithInjectNow() Option {
	return func(p *Processor) {
		p.injectNow = true
	}
}

// New creates a Processor that operates on the given field name.
// The field is created when absent; existing values are re-formatted
// unless WithInjectNow is specified.
func New(field string, opts ...Option) (*Processor, error) {
	if field == "" {
		return nil, fmt.Errorf("timestamp: field name must not be empty")
	}
	p := &Processor{
		field:  field,
		format: time.RFC3339,
	}
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// Apply processes a single JSON log line and returns the modified line.
func (p *Processor) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("timestamp: invalid JSON: %w", err)
	}

	if p.injectNow {
		m[p.field] = time.Now().UTC().Format(p.format)
	} else {
		switch v := m[p.field].(type) {
		case string:
			t, err := time.Parse(p.format, v)
			if err != nil {
				// leave the field untouched if it cannot be parsed
				break
			}
			m[p.field] = t.UTC().Format(p.format)
		case nil:
			m[p.field] = time.Now().UTC().Format(p.format)
		}
	}

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("timestamp: marshal error: %w", err)
	}
	return out, nil
}
