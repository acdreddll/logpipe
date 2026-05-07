package fieldround

import (
	"encoding/json"
	"fmt"
	"math"
)

// Rounder rounds a numeric JSON field to a given number of decimal places.
type Rounder struct {
	field  string
	places int
}

// New returns a Rounder that rounds the value at field to places decimal places.
// places must be >= 0.
func New(field string, places int) (*Rounder, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldround: field name must not be empty")
	}
	if places < 0 {
		return nil, fmt.Errorf("fieldround: places must be >= 0, got %d", places)
	}
	return &Rounder{field: field, places: places}, nil
}

// Apply rounds the target field in the JSON line and returns the modified line.
// If the field is absent the line is returned unchanged.
// Returns an error if line is not valid JSON or the field is not numeric.
func (r *Rounder) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("fieldround: invalid JSON: %w", err)
	}

	v, ok := m[r.field]
	if !ok {
		return line, nil
	}

	f, ok := toFloat(v)
	if !ok {
		return nil, fmt.Errorf("fieldround: field %q is not numeric", r.field)
	}

	mul := math.Pow(10, float64(r.places))
	m[r.field] = math.Round(f*mul) / mul

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("fieldround: marshal: %w", err)
	}
	return out, nil
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}
