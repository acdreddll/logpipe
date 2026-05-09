// Package fieldabs provides a processor that replaces a numeric JSON field
// with its absolute value.
package fieldabs

import (
	"encoding/json"
	"fmt"
	"math"
)

// Processor replaces a numeric field with its absolute value.
type Processor struct {
	field string
}

// New returns a Processor that operates on the given field.
// An error is returned if field is empty.
func New(field string) (*Processor, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldabs: field name must not be empty")
	}
	return &Processor{field: field}, nil
}

// Apply reads the named field from the JSON line, replaces it with its
// absolute value, and returns the updated JSON. Non-numeric fields and
// missing fields are left unchanged. Invalid JSON returns an error.
func (p *Processor) Apply(line string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldabs: invalid JSON: %w", err)
	}

	v, ok := m[p.field]
	if !ok {
		return line, nil
	}

	switch n := v.(type) {
	case float64:
		m[p.field] = math.Abs(n)
	default:
		return line, nil
	}

	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldabs: marshal: %w", err)
	}
	return string(out), nil
}
