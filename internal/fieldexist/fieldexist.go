// Package fieldexist provides a processor that filters log events based on
// the presence or absence of one or more JSON fields.
package fieldexist

import (
	"encoding/json"
	"fmt"
)

// Mode controls whether events are kept when fields exist or when they are absent.
type Mode int

const (
	// ModeRequire keeps events that have ALL of the specified fields.
	ModeRequire Mode = iota
	// ModeExclude keeps events that have NONE of the specified fields.
	ModeExclude
)

// Filter keeps or drops log events based on field presence.
type Filter struct {
	fields []string
	mode   Mode
}

// New creates a Filter that operates in the given mode over the provided fields.
// At least one field name must be supplied and no field name may be empty.
func New(mode Mode, fields []string) (*Filter, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("fieldexist: at least one field is required")
	}
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("fieldexist: field name must not be empty")
		}
	}
	if mode != ModeRequire && mode != ModeExclude {
		return nil, fmt.Errorf("fieldexist: unknown mode %d", mode)
	}
	copy := make([]string, len(fields))
	copy = append(copy[:0], fields...)
	return &Filter{fields: copy, mode: mode}, nil
}

// Keep returns true when the event should be forwarded downstream.
// It returns an error if line is not valid JSON.
func (f *Filter) Keep(line []byte) (bool, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(line, &m); err != nil {
		return false, fmt.Errorf("fieldexist: invalid JSON: %w", err)
	}

	switch f.mode {
	case ModeRequire:
		for _, field := range f.fields {
			if _, ok := m[field]; !ok {
				return false, nil
			}
		}
		return true, nil
	case ModeExclude:
		for _, field := range f.fields {
			if _, ok := m[field]; ok {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}
