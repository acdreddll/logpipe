// Package transform provides log record transformation capabilities.
package transform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Op represents a transformation operation type.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpRename Op = "rename"
)

// Rule defines a single transformation rule.
type Rule struct {
	Op    Op     `json:"op"`
	Field string `json:"field"`
	Value string `json:"value,omitempty"`
	To    string `json:"to,omitempty"`
}

// Transformer applies a set of rules to log records.
type Transformer struct {
	rules []Rule
}

// New creates a Transformer from a slice of rules.
func New(rules []Rule) (*Transformer, error) {
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("rule %d: field is required", i)
		}
		switch r.Op {
		case OpSet, OpDelete, OpRename:
			// valid op
		default:
			return nil, fmt.Errorf("rule %d: unknown op %q", i, r.Op)
		}
		if r.Op == OpRename && r.To == "" {
			return nil, fmt.Errorf("rule %d: 'to' is required for rename op", i)
		}
	}
	return &Transformer{rules: rules}, nil
}

// Apply transforms a raw JSON log line and returns the modified JSON.
func (t *Transformer) Apply(line string) (string, error) {
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &record); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	for _, r := range t.rules {
		switch r.Op {
		case OpSet:
			record[r.Field] = r.Value
		case OpDelete:
			delete(record, r.Field)
		case OpRename:
			if v, ok := record[r.Field]; ok {
				record[r.To] = v
				delete(record, r.Field)
			}
		}
	}
	out, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Rules returns a copy of the rules held by the Transformer.
func (t *Transformer) Rules() []Rule {
	copy := make([]Rule, len(t.rules))
	for i, r := range t.rules {
		copy[i] = r
	}
	return copy
}
