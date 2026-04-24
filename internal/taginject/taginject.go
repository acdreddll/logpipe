// Package taginject provides a processor that injects a static or
// computed tag field into every log line passing through the pipeline.
package taginject

import (
	"encoding/json"
	"fmt"
)

// Injector adds a fixed tag field to each JSON log line.
type Injector struct {
	field string
	value string
}

// Option is a functional option for Injector.
type Option func(*Injector)

// WithField overrides the default field name ("tag").
func WithField(field string) Option {
	return func(i *Injector) {
		i.field = field
	}
}

// New creates an Injector that sets field to value on every log line.
// value must be a non-empty string.
func New(value string, opts ...Option) (*Injector, error) {
	if value == "" {
		return nil, fmt.Errorf("taginject: value must not be empty")
	}
	i := &Injector{
		field: "tag",
		value: value,
	}
	for _, o := range opts {
		o(i)
	}
	if i.field == "" {
		return nil, fmt.Errorf("taginject: field must not be empty")
	}
	return i, nil
}

// Apply injects the tag into the JSON line and returns the modified line.
// If the field already exists it is overwritten.
func (i *Injector) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("taginject: invalid JSON: %w", err)
	}
	m[i.field] = i.value
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("taginject: marshal error: %w", err)
	}
	return out, nil
}
