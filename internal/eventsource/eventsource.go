// Package eventsource provides a named source tag injector that stamps
// each log line with a configurable source identifier (e.g. "app", "nginx").
// It is useful when merging streams from multiple origins so that downstream
// routes can filter by origin.
package eventsource

import (
	"encoding/json"
	"fmt"
)

const defaultField = "source"

// Source stamps log lines with a fixed source name.
type Source struct {
	field string
	name  string
}

// Option is a functional option for Source.
type Option func(*Source)

// WithField overrides the JSON field name (default: "source").
func WithField(f string) Option {
	return func(s *Source) { s.field = f }
}

// New creates a Source that injects name into every log line.
// name must be non-empty. An optional field name can be supplied via WithField.
func New(name string, opts ...Option) (*Source, error) {
	if name == "" {
		return nil, fmt.Errorf("eventsource: name must not be empty")
	}
	s := &Source{field: defaultField, name: name}
	for _, o := range opts {
		o(s)
	}
	if s.field == "" {
		return nil, fmt.Errorf("eventsource: field must not be empty")
	}
	return s, nil
}

// Apply injects the source field into the JSON log line.
// If the field already exists it is overwritten.
func (s *Source) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("eventsource: invalid JSON: %w", err)
	}
	m[s.field] = s.name
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("eventsource: marshal error: %w", err)
	}
	return out, nil
}

// Name returns the source name this stamper injects.
func (s *Source) Name() string { return s.name }

// Field returns the JSON field name used for injection.
func (s *Source) Field() string { return s.field }
