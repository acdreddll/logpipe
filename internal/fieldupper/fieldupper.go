// Package fieldupper converts the string value of a JSON field to uppercase.
package fieldupper

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Transformer uppercases the value of a single JSON field.
type Transformer struct {
	field string
}

// New returns a Transformer that uppercases the value of field.
// An error is returned if field is empty.
func New(field string) (*Transformer, error) {
	if strings.TrimSpace(field) == "" {
		return nil, fmt.Errorf("fieldupper: field name must not be empty")
	}
	return &Transformer{field: field}, nil
}

// Apply reads line as a JSON object, uppercases the target field value,
// and returns the modified JSON. If the field is absent or not a string
// the line is returned unchanged. Non-JSON input returns an error.
func (t *Transformer) Apply(line string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldupper: invalid JSON: %w", err)
	}

	v, ok := m[t.field]
	if !ok {
		return line, nil
	}

	s, ok := v.(string)
	if !ok {
		return line, nil
	}

	m[t.field] = strings.ToUpper(s)

	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldupper: marshal: %w", err)
	}
	return string(out), nil
}
