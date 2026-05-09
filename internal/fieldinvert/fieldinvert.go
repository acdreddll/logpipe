// Package fieldinvert provides a processor that inverts the boolean value
// of a named field within a JSON log event.
package fieldinvert

import (
	"encoding/json"
	"fmt"
)

// Inverter inverts the boolean value of a specific field in a JSON log line.
type Inverter struct {
	field string
}

// New returns a new Inverter that will invert the boolean at the given field.
// Returns an error if field is empty.
func New(field string) (*Inverter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldinvert: field name must not be empty")
	}
	return &Inverter{field: field}, nil
}

// Apply inverts the boolean value of the configured field in the JSON line.
// If the field is absent the line is returned unchanged.
// Returns an error if the line is not valid JSON or the field is not a boolean.
func (inv *Inverter) Apply(line string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldinvert: invalid JSON: %w", err)
	}

	val, ok := m[inv.field]
	if !ok {
		return line, nil
	}

	b, ok := val.(bool)
	if !ok {
		return "", fmt.Errorf("fieldinvert: field %q is not a boolean", inv.field)
	}

	m[inv.field] = !b

	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldinvert: marshal: %w", err)
	}
	return string(out), nil
}
