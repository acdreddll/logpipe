// Package fieldboolean coerces a JSON field value to a boolean.
//
// Truthy string values: "true", "1", "yes", "on"  (case-insensitive).
// Falsy string values : "false", "0", "no", "off" (case-insensitive).
// Numeric values      : non-zero → true, zero → false.
// Existing booleans are left unchanged.
package fieldboolean

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Converter coerces a single field to a boolean value.
type Converter struct {
	field string
}

// New returns a Converter that coerces the named field to a boolean.
// field must not be empty.
func New(field string) (*Converter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldboolean: field name must not be empty")
	}
	return &Converter{field: field}, nil
}

// Apply coerces the target field in the JSON-encoded line and returns the
// updated JSON. If the field is absent the line is returned unchanged.
func (c *Converter) Apply(line []byte) ([]byte, error) {
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("fieldboolean: invalid JSON: %w", err)
	}

	v, ok := m[c.field]
	if !ok {
		return line, nil
	}

	b, err := toBool(v)
	if err != nil {
		return nil, fmt.Errorf("fieldboolean: field %q: %w", c.field, err)
	}
	m[c.field] = b

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("fieldboolean: marshal: %w", err)
	}
	return out, nil
}

func toBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "true", "1", "yes", "on":
			return true, nil
		case "false", "0", "no", "off", "":
			return false, nil
		default:
			return false, fmt.Errorf("unrecognised boolean string %q", val)
		}
	case float64:
		return val != 0, nil
	default:
		return false, fmt.Errorf("unsupported type %T", v)
	}
}
