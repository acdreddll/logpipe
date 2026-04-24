// Package labelfilter provides a processor that keeps or removes log lines
// based on whether a specified field's value matches a set of allowed or
// denied labels.
package labelfilter

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Mode controls whether the label list is an allow-list or a deny-list.
type Mode string

const (
	ModeAllow Mode = "allow"
	ModeDeny  Mode = "deny"
)

// LabelFilter decides whether a log line should pass through based on a
// field's value and a set of labels.
type LabelFilter struct {
	field  string
	mode   Mode
	labels map[string]struct{}
}

// New creates a LabelFilter that operates on field using mode and the
// provided label set. At least one label must be supplied.
func New(field string, mode Mode, labels []string) (*LabelFilter, error) {
	if field == "" {
		return nil, errors.New("labelfilter: field must not be empty")
	}
	if mode != ModeAllow && mode != ModeDeny {
		return nil, fmt.Errorf("labelfilter: unknown mode %q", mode)
	}
	if len(labels) == 0 {
		return nil, errors.New("labelfilter: at least one label is required")
	}
	set := make(map[string]struct{}, len(labels))
	for _, l := range labels {
		set[l] = struct{}{}
	}
	return &LabelFilter{field: field, mode: mode, labels: set}, nil
}

// Keep returns true when the log line should be forwarded downstream.
// Lines that are not valid JSON are dropped (returns false, nil).
func (lf *LabelFilter) Keep(line []byte) (bool, error) {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return false, nil
	}

	v, ok := obj[lf.field]
	if !ok {
		// Field absent: allow-list drops it, deny-list keeps it.
		return lf.mode == ModeDeny, nil
	}

	val := fmt.Sprintf("%v", v)
	_, matched := lf.labels[val]

	switch lf.mode {
	case ModeAllow:
		return matched, nil
	case ModeDeny:
		return !matched, nil
	}
	return true, nil
}
