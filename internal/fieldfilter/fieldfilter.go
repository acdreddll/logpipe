// Package fieldfilter provides a processor that retains or removes a
// specific set of fields from a JSON log line.
package fieldfilter

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Mode controls whether the listed fields are kept or dropped.
type Mode int

const (
	// ModeAllow retains only the listed fields.
	ModeAllow Mode = iota
	// ModeDeny removes the listed fields.
	ModeDeny
)

// FieldFilter selectively keeps or removes fields from a JSON log line.
type FieldFilter struct {
	fields map[string]struct{}
	mode   Mode
}

// New creates a FieldFilter operating in the given mode over the
// provided field names. At least one field name must be supplied.
func New(mode Mode, fields []string) (*FieldFilter, error) {
	if len(fields) == 0 {
		return nil, errors.New("fieldfilter: at least one field name is required")
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if f == "" {
			return nil, errors.New("fieldfilter: field name must not be empty")
		}
		set[f] = struct{}{}
	}
	return &FieldFilter{fields: set, mode: mode}, nil
}

// Apply processes a single JSON log line, returning a new line with
// fields retained or removed according to the configured mode.
func (ff *FieldFilter) Apply(line []byte) ([]byte, error) {
	var record map[string]json.RawMessage
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("fieldfilter: invalid JSON: %w", err)
	}

	switch ff.mode {
	case ModeAllow:
		filtered := make(map[string]json.RawMessage, len(ff.fields))
		for k, v := range record {
			if _, ok := ff.fields[k]; ok {
				filtered[k] = v
			}
		}
		record = filtered
	case ModeDeny:
		for k := range ff.fields {
			delete(record, k)
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("fieldfilter: marshal error: %w", err)
	}
	return out, nil
}
