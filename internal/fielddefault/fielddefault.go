// Package fielddefault sets a default value on a JSON log field when the
// field is absent or empty.
package fielddefault

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Defaulter sets a default value on a named field.
type Defaulter struct {
	field    string
	value    any
	overwrite bool
}

// Option is a functional option for Defaulter.
type Option func(*Defaulter)

// WithOverwrite causes the defaulter to overwrite the field even when it
// already has a value.
func WithOverwrite() Option {
	return func(d *Defaulter) { d.overwrite = true }
}

// New creates a Defaulter that injects value into field when the field is
// absent (or empty string). Pass WithOverwrite to always set the field.
func New(field string, value any, opts ...Option) (*Defaulter, error) {
	if field == "" {
		return nil, errors.New("fielddefault: field name must not be empty")
	}
	if value == nil {
		return nil, errors.New("fielddefault: value must not be nil")
	}
	d := &Defaulter{field: field, value: value}
	for _, o := range opts {
		o(d)
	}
	return d, nil
}

// Apply injects the default value into the JSON log line and returns the
// modified line. The input must be a valid JSON object.
func (d *Defaulter) Apply(line []byte) ([]byte, error) {
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("fielddefault: invalid JSON: %w", err)
	}

	existing, exists := m[d.field]
	if !exists || existing == nil || existing == "" || d.overwrite {
		m[d.field] = d.value
	}

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("fielddefault: marshal error: %w", err)
	}
	return out, nil
}
