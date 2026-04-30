// Package fieldsplit splits a string field into a JSON array using a delimiter.
package fieldsplit

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Splitter splits a named string field into a JSON array.
type Splitter struct {
	field string
	dest  string
	sep   string
}

// Option is a functional option for Splitter.
type Option func(*Splitter)

// WithDest overrides the destination field (defaults to the source field).
func WithDest(dest string) Option {
	return func(s *Splitter) { s.dest = dest }
}

// WithSeparator sets the delimiter (defaults to ",").
func WithSeparator(sep string) Option {
	return func(s *Splitter) { s.sep = sep }
}

// New creates a Splitter for the given field.
func New(field string, opts ...Option) (*Splitter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldsplit: field name must not be empty")
	}
	s := &Splitter{field: field, dest: field, sep: ","}
	for _, o := range opts {
		o(s)
	}
	if s.dest == "" {
		return nil, fmt.Errorf("fieldsplit: dest field name must not be empty")
	}
	if s.sep == "" {
		return nil, fmt.Errorf("fieldsplit: separator must not be empty")
	}
	return s, nil
}

// Apply parses line as JSON, splits the configured field, and returns the
// modified JSON. If the field is absent the line is returned unchanged.
func (s *Splitter) Apply(line string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldsplit: invalid JSON: %w", err)
	}

	v, ok := m[s.field]
	if !ok {
		return line, nil
	}

	str, ok := v.(string)
	if !ok {
		return line, nil
	}

	parts := strings.Split(str, s.sep)
	result := make([]any, len(parts))
	for i, p := range parts {
		result[i] = strings.TrimSpace(p)
	}
	m[s.dest] = result

	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldsplit: marshal: %w", err)
	}
	return string(b), nil
}
