// Package fieldcopy copies the value of one JSON field into another field,
// leaving the source field intact.
package fieldcopy

import (
	"encoding/json"
	"fmt"
)

// Copier copies a source field value into a destination field.
type Copier struct {
	src string
	dst string
	overwrite bool
}

// Option is a functional option for Copier.
type Option func(*Copier)

// WithOverwrite allows the destination field to be overwritten if it already exists.
func WithOverwrite() Option {
	return func(c *Copier) { c.overwrite = true }
}

// New creates a new Copier that copies src into dst.
func New(src, dst string, opts ...Option) (*Copier, error) {
	if src == "" {
		return nil, fmt.Errorf("fieldcopy: source field must not be empty")
	}
	if dst == "" {
		return nil, fmt.Errorf("fieldcopy: destination field must not be empty")
	}
	c := &Copier{src: src, dst: dst}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// Apply copies the source field value into the destination field on the JSON
// log line. The original line is returned unchanged on any error.
func (c *Copier) Apply(line string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line, fmt.Errorf("fieldcopy: invalid JSON: %w", err)
	}

	val, ok := m[c.src]
	if !ok {
		return line, nil
	}

	if _, exists := m[c.dst]; exists && !c.overwrite {
		return line, nil
	}

	m[c.dst] = val

	out, err := json.Marshal(m)
	if err != nil {
		return line, fmt.Errorf("fieldcopy: marshal error: %w", err)
	}
	return string(out), nil
}
