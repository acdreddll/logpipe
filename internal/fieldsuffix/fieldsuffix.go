// Package fieldsuffix appends a static suffix string to a named string field
// within a JSON log event.
package fieldsuffix

import (
	"encoding/json"
	"fmt"
)

// Suffixer appends a fixed string to a JSON field value.
type Suffixer struct {
	field  string
	suffix string
}

// New returns a Suffixer that appends suffix to the value of field.
// Both field and suffix must be non-empty.
func New(field, suffix string) (*Suffixer, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldsuffix: field must not be empty")
	}
	if suffix == "" {
		return nil, fmt.Errorf("fieldsuffix: suffix must not be empty")
	}
	return &Suffixer{field: field, suffix: suffix}, nil
}

// Apply appends the configured suffix to the named field in the JSON line.
// If the field is absent or not a string the line is returned unchanged.
// Invalid JSON returns an error.
func (s *Suffixer) Apply(line string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldsuffix: invalid JSON: %w", err)
	}

	v, ok := m[s.field]
	if !ok {
		return line, nil
	}
	str, ok := v.(string)
	if !ok {
		return line, nil
	}

	m[s.field] = str + s.suffix

	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldsuffix: marshal: %w", err)
	}
	return string(out), nil
}
