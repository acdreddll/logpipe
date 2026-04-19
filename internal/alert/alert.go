package alert

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Condition defines when an alert should fire.
type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // eq, contains, exists
	Value    string `json:"value,omitempty"`
}

// Alert fires a notification when a log line matches a condition.
type Alert struct {
	Name      string
	condition Condition
	notify    func(name, line string)
}

// New creates an Alert with the given condition and notify callback.
func New(name string, c Condition, notify func(name, line string)) (*Alert, error) {
	if name == "" {
		return nil, fmt.Errorf("alert name must not be empty")
	}
	if c.Field == "" {
		return nil, fmt.Errorf("alert condition field must not be empty")
	}
	validOps := map[string]bool{"eq": true, "contains": true, "exists": true}
	if !validOps[c.Operator] {
		return nil, fmt.Errorf("unknown operator %q", c.Operator)
	}
	if notify == nil {
		return nil, fmt.Errorf("notify callback must not be nil")
	}
	return &Alert{Name: name, condition: c, notify: notify}, nil
}

// Evaluate checks the log line against the condition and fires notify if matched.
func (a *Alert) Evaluate(line string) error {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	val, exists := m[a.condition.Field]
	switch a.condition.Operator {
	case "exists":
		if exists {
			a.notify(a.Name, line)
		}
	case "eq":
		if exists && fmt.Sprintf("%v", val) == a.condition.Value {
			a.notify(a.Name, line)
		}
	case "contains":
		if exists && strings.Contains(fmt.Sprintf("%v", val), a.condition.Value) {
			a.notify(a.Name, line)
		}
	}
	return nil
}
