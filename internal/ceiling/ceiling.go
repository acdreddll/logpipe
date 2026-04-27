// Package ceiling provides a processor that clamps numeric field values
// to a configurable maximum (ceiling). Values exceeding the ceiling are
// replaced with the ceiling value; non-numeric or absent fields are left
// unchanged.
package ceiling

import (
	"encoding/json"
	"fmt"
	"math"
)

// Ceiling clamps a named numeric field to a maximum value.
type Ceiling struct {
	field string
	max   float64
}

// New returns a Ceiling that clamps field to max.
// field must be non-empty and max must be finite.
func New(field string, max float64) (*Ceiling, error) {
	if field == "" {
		return nil, fmt.Errorf("ceiling: field name must not be empty")
	}
	if math.IsInf(max, 0) || math.IsNaN(max) {
		return nil, fmt.Errorf("ceiling: max must be a finite number")
	}
	return &Ceiling{field: field, max: max}, nil
}

// Apply clamps the numeric value at c.field in the JSON log line.
// Lines with a missing or non-numeric field are returned unchanged.
// Invalid JSON returns an error.
func (c *Ceiling) Apply(line []byte) ([]byte, error) {
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("ceiling: invalid JSON: %w", err)
	}

	v, ok := m[c.field]
	if !ok {
		return line, nil
	}

	f, ok := toFloat(v)
	if !ok {
		return line, nil
	}

	if f > c.max {
		m[c.field] = c.max
		out, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("ceiling: marshal: %w", err)
		}
		return out, nil
	}

	return line, nil
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}
