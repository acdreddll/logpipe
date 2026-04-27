// Package fieldclone copies the value of one JSON field into one or more
// destination fields, leaving the source field intact.
package fieldclone

import (
	"encoding/json"
	"fmt"
)

// Cloner copies a source field value into one or more destination fields.
type Cloner struct {
	src  string
	dsts []string
}

// New creates a Cloner that copies src into each field listed in dsts.
// Returns an error if src is empty or dsts is empty.
func New(src string, dsts []string) (*Cloner, error) {
	if src == "" {
		return nil, fmt.Errorf("fieldclone: source field must not be empty")
	}
	if len(dsts) == 0 {
		return nil, fmt.Errorf("fieldclone: at least one destination field is required")
	}
	for i, d := range dsts {
		if d == "" {
			return nil, fmt.Errorf("fieldclone: destination field at index %d must not be empty", i)
		}
	}
	return &Cloner{src: src, dsts: dsts}, nil
}

// Apply clones the source field value into every destination field.
// If the source field is absent the line is returned unchanged.
// Returns an error if line is not valid JSON.
func (c *Cloner) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("fieldclone: invalid JSON: %w", err)
	}

	val, ok := m[c.src]
	if !ok {
		return line, nil
	}

	for _, d := range c.dsts {
		m[d] = val
	}

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("fieldclone: marshal error: %w", err)
	}
	return out, nil
}
