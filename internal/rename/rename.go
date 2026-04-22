// Package rename provides a processor that renames a set of JSON log fields
// according to a static mapping, leaving all other fields untouched.
package rename

import (
	"encoding/json"
	"fmt"
)

// Renamer renames fields in a JSON log line according to a fixed mapping.
type Renamer struct {
	// mapping is old-name → new-name.
	mapping map[string]string
}

// New creates a Renamer from the provided mapping.
// mapping must not be empty and neither key nor value may be the empty string.
func New(mapping map[string]string) (*Renamer, error) {
	if len(mapping) == 0 {
		return nil, fmt.Errorf("rename: mapping must not be empty")
	}
	for k, v := range mapping {
		if k == "" {
			return nil, fmt.Errorf("rename: source field name must not be empty")
		}
		if v == "" {
			return nil, fmt.Errorf("rename: destination field name must not be empty (source: %q)", k)
		}
	}
	copy := make(map[string]string, len(mapping))
	for k, v := range mapping {
		copy[k] = v
	}
	return &Renamer{mapping: copy}, nil
}

// Apply renames fields in line according to the configured mapping.
// Fields not present in the mapping are passed through unchanged.
// If a destination field already exists it is overwritten.
// Returns an error when line is not valid JSON.
func (r *Renamer) Apply(line []byte) ([]byte, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("rename: invalid JSON: %w", err)
	}

	for src, dst := range r.mapping {
		val, ok := obj[src]
		if !ok {
			continue
		}
		delete(obj, src)
		obj[dst] = val
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("rename: marshal: %w", err)
	}
	return out, nil
}
