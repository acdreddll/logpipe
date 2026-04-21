// Package fieldmap provides a Mapper that copies or moves log fields
// from one key to another, optionally removing the source key.
package fieldmap

import (
	"encoding/json"
	"fmt"
)

// Mapping describes a single field remapping rule.
type Mapping struct {
	From   string
	To     string
	Delete bool // if true, remove the source field after copying
}

// Mapper applies a set of field mappings to JSON log lines.
type Mapper struct {
	mappings []Mapping
}

// New creates a Mapper from the provided mappings.
// Returns an error if any mapping has an empty From or To field.
func New(mappings []Mapping) (*Mapper, error) {
	for i, m := range mappings {
		if m.From == "" {
			return nil, fmt.Errorf("mapping[%d]: From field must not be empty", i)
		}
		if m.To == "" {
			return nil, fmt.Errorf("mapping[%d]: To field must not be empty", i)
		}
	}
	return &Mapper{mappings: mappings}, nil
}

// Apply processes a single JSON log line and returns the remapped line.
// Fields that do not exist in the source are silently skipped.
func (m *Mapper) Apply(line []byte) ([]byte, error) {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("fieldmap: invalid JSON: %w", err)
	}

	for _, mapping := range m.mappings {
		val, ok := obj[mapping.From]
		if !ok {
			continue
		}
		obj[mapping.To] = val
		if mapping.Delete {
			delete(obj, mapping.From)
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("fieldmap: marshal: %w", err)
	}
	return out, nil
}
