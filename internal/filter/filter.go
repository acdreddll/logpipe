package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a single filter rule applied to a log entry field.
type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // eq, contains, exists, neq
	Value    string `json:"value"`
}

// Filter holds a set of rules and applies them to log entries.
type Filter struct {
	Rules []Rule
}

// New creates a Filter from the given rules.
func New(rules []Rule) *Filter {
	return &Filter{Rules: rules}
}

// Match returns true if the log line (JSON) satisfies all rules.
func (f *Filter) Match(line []byte) (bool, error) {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return false, fmt.Errorf("filter: failed to parse log line: %w", err)
	}

	for _, rule := range f.Rules {
		val, exists := entry[rule.Field]
		switch rule.Operator {
		case "exists":
			if !exists {
				return false, nil
			}
		case "eq":
			if !exists || toString(val) != rule.Value {
				return false, nil
			}
		case "neq":
			if exists && toString(val) == rule.Value {
				return false, nil
			}
		case "contains":
			if !exists || !strings.Contains(toString(val), rule.Value) {
				return false, nil
			}
		default:
			return false, fmt.Errorf("filter: unknown operator %q", rule.Operator)
		}
	}
	return true, nil
}

func toString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
