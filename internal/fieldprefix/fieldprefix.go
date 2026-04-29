// Package fieldprefix adds a string prefix to the value of a named JSON field.
package fieldprefix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Prefixer prepends a fixed string to a field value.
type Prefixer struct {
	field  string
	prefix string
}

// New returns a Prefixer that prepends prefix to the value of field.
// Both field and prefix must be non-empty.
func New(field, prefix string) (*Prefixer, error) {
	if field == "" {
		return nil, errors.New("fieldprefix: field must not be empty")
	}
	if prefix == "" {
		return nil, errors.New("fieldprefix: prefix must not be empty")
	}
	return &Prefixer{field: field, prefix: prefix}, nil
}

// Apply prepends the configured prefix to the target field in the JSON line.
// If the field is absent or not a string the line is returned unchanged.
// An error is returned when line is not valid JSON.
func (p *Prefixer) Apply(line []byte) ([]byte, error) {
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("fieldprefix: invalid JSON: %w", err)
	}

	v, ok := m[p.field]
	if !ok {
		return line, nil
	}

	s, ok := v.(string)
	if !ok {
		return line, nil
	}

	m[p.field] = p.prefix + s

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("fieldprefix: marshal: %w", err)
	}
	return out, nil
}
