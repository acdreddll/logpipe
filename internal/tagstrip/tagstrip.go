// Package tagstrip removes specified tag fields from structured log lines.
package tagstrip

import (
	"encoding/json"
	"fmt"
)

// Stripper removes a set of named fields from a JSON log line.
type Stripper struct {
	fields map[string]struct{}
}

// New creates a Stripper that removes the given field names.
// Returns an error if fields is empty or any name is blank.
func New(fields []string) (*Stripper, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("tagstrip: at least one field name is required")
	}
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("tagstrip: field name must not be empty")
		}
		m[f] = struct{}{}
	}
	return &Stripper{fields: m}, nil
}

// Apply removes the configured fields from the JSON log line.
// Returns the modified line, or an error if the input is not valid JSON.
func (s *Stripper) Apply(line []byte) ([]byte, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("tagstrip: invalid JSON: %w", err)
	}
	for f := range s.fields {
		delete(obj, f)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("tagstrip: marshal error: %w", err)
	}
	return out, nil
}

// Fields returns the set of field names that will be stripped.
func (s *Stripper) Fields() []string {
	out := make([]string, 0, len(s.fields))
	for f := range s.fields {
		out = append(out, f)
	}
	return out
}
