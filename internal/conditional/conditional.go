// Package conditional provides a processor that applies a transformation
// only when a specified field matches a given value.
package conditional

import (
	"encoding/json"
	"fmt"
)

// Processor applies an inner transform only when a condition is met.
type Processor struct {
	field    string
	operator string
	value    string
	apply    func([]byte) ([]byte, error)
}

// New creates a Processor that calls apply when field op value is true.
// Supported operators: "eq", "neq", "exists".
func New(field, operator, value string, apply func([]byte) ([]byte, error)) (*Processor, error) {
	if field == "" {
		return nil, fmt.Errorf("conditional: field must not be empty")
	}
	switch operator {
	case "eq", "neq", "exists":
	default:
		return nil, fmt.Errorf("conditional: unsupported operator %q", operator)
	}
	if apply == nil {
		return nil, fmt.Errorf("conditional: apply func must not be nil")
	}
	return &Processor{field: field, operator: operator, value: value, apply: apply}, nil
}

// Apply evaluates the condition against line and, if true, runs the inner
// transform. If the condition is false the original line is returned unchanged.
func (p *Processor) Apply(line []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("conditional: invalid JSON: %w", err)
	}

	raw, exists := obj[p.field]

	var matched bool
	switch p.operator {
	case "exists":
		matched = exists
	case "eq":
		matched = exists && fmt.Sprintf("%v", raw) == p.value
	case "neq":
		matched = !exists || fmt.Sprintf("%v", raw) != p.value
	}

	if !matched {
		return line, nil
	}
	return p.apply(line)
}
