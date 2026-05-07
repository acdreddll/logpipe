// Package fieldtrim trims leading and trailing whitespace (or a custom
// cutset) from a string field in a JSON log event.
package fieldtrim

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Trimmer trims characters from a named field value.
type Trimmer struct {
	field  string
	cutset string // empty → trim whitespace
}

// Option is a functional option for Trimmer.
type Option func(*Trimmer)

// WithCutset sets an explicit cutset of characters to trim instead of
// Unicode whitespace.
func WithCutset(cutset string) Option {
	return func(t *Trimmer) { t.cutset = cutset }
}

// New creates a Trimmer that operates on field. At least one character in
// the field value must be non-whitespace (or non-cutset) for a change to
// occur; missing fields are silently ignored.
func New(field string, opts ...Option) (*Trimmer, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldtrim: field name must not be empty")
	}
	t := &Trimmer{field: field}
	for _, o := range opts {
		o(t)
	}
	return t, nil
}

// Apply trims the target field in the JSON-encoded log line and returns the
// modified line. Lines that are not valid JSON or that do not contain the
// field are returned unchanged.
func (t *Trimmer) Apply(line string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldtrim: invalid JSON: %w", err)
	}

	v, ok := m[t.field]
	if !ok {
		return line, nil
	}

	s, ok := v.(string)
	if !ok {
		return line, nil
	}

	var trimmed string
	if t.cutset == "" {
		trimmed = strings.TrimSpace(s)
	} else {
		trimmed = strings.Trim(s, t.cutset)
	}

	if trimmed == s {
		return line, nil
	}

	m[t.field] = trimmed
	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldtrim: marshal error: %w", err)
	}
	return string(out), nil
}
