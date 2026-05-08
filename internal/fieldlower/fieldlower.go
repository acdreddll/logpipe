package fieldlower

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Lowercaser lowercases the string value of a single JSON field.
type Lowercaser struct {
	field string
}

// New returns a Lowercaser that lowercases the value of the given field.
// An error is returned if field is empty.
func New(field string) (*Lowercaser, error) {
	if strings.TrimSpace(field) == "" {
		return nil, fmt.Errorf("fieldlower: field name must not be empty")
	}
	return &Lowercaser{field: field}, nil
}

// Apply lowercases the target field in the JSON-encoded log line.
// Non-string fields and absent fields are left unchanged.
// The original line is returned unmodified on any parse error.
func (l *Lowercaser) Apply(line string) string {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}

	v, ok := m[l.field]
	if !ok {
		return line
	}

	s, ok := v.(string)
	if !ok {
		return line
	}

	m[l.field] = strings.ToLower(s)

	out, err := json.Marshal(m)
	if err != nil {
		return line
	}
	return string(out)
}
