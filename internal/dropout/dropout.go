// Package dropout provides a log-line dropper that discards lines
// whose JSON field value matches a set of configured literal strings.
// It is useful for suppressing known-noisy or low-value log entries
// before they reach downstream outputs.
package dropout

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Dropper discards log lines where a named field equals one of a set
// of configured values.
type Dropper struct {
	field  string
	values map[string]struct{}
}

// New creates a Dropper that drops any line whose JSON field named
// field equals one of the provided values.
// field must be non-empty and values must contain at least one entry.
func New(field string, values []string) (*Dropper, error) {
	if field == "" {
		return nil, errors.New("dropout: field must not be empty")
	}
	if len(values) == 0 {
		return nil, errors.New("dropout: at least one value is required")
	}
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return &Dropper{field: field, values: set}, nil
}

// ShouldDrop returns true when line is valid JSON and the configured
// field's value is present in the drop-set. Malformed JSON lines are
// never dropped (returns false, non-nil error).
func (d *Dropper) ShouldDrop(line []byte) (bool, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return false, fmt.Errorf("dropout: invalid JSON: %w", err)
	}
	v, ok := obj[d.field]
	if !ok {
		return false, nil
	}
	s, ok := v.(string)
	if !ok {
		return false, nil
	}
	_, drop := d.values[s]
	return drop, nil
}
