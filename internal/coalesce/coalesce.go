// Package coalesce provides a processor that returns the first non-empty
// value from a list of fields and writes it to a destination field.
package coalesce

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Coalescer reads a prioritised list of source fields and writes the first
// non-empty string value it finds into a destination field.
type Coalescer struct {
	sources []string
	dest    string
}

// New creates a Coalescer that searches sources in order and writes the first
// non-empty value to dest. At least one source and a non-empty dest are
// required.
func New(dest string, sources []string) (*Coalescer, error) {
	if dest == "" {
		return nil, errors.New("coalesce: dest field must not be empty")
	}
	if len(sources) == 0 {
		return nil, errors.New("coalesce: at least one source field is required")
	}
	return &Coalescer{dest: dest, sources: sources}, nil
}

// Apply reads line as a JSON object, finds the first non-empty source field,
// writes its value to dest, and returns the modified JSON. If no source field
// has a non-empty value the line is returned unchanged. An error is returned
// for invalid JSON.
func (c *Coalescer) Apply(line string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("coalesce: invalid JSON: %w", err)
	}

	for _, src := range c.sources {
		v, ok := m[src]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		m[c.dest] = s
		out, err := json.Marshal(m)
		if err != nil {
			return "", fmt.Errorf("coalesce: marshal error: %w", err)
		}
		return string(out), nil
	}

	return line, nil
}
