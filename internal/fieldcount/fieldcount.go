// Package fieldcount injects the number of keys present in a JSON log event
// into a configurable destination field.
package fieldcount

import (
	"encoding/json"
	"fmt"
)

// defaultField is the destination field used when none is specified.
const defaultField = "_field_count"

// Counter counts the number of top-level fields in a JSON event and injects
// the result into a target field.
type Counter struct {
	field string
}

// WithField sets the destination field name.
func WithField(field string) func(*Counter) {
	return func(c *Counter) {
		c.field = field
	}
}

// New creates a Counter. An empty field name is rejected.
func New(opts ...func(*Counter)) (*Counter, error) {
	c := &Counter{field: defaultField}
	for _, o := range opts {
		o(c)
	}
	if c.field == "" {
		return nil, fmt.Errorf("fieldcount: field name must not be empty")
	}
	return c, nil
}

// Apply parses line as a JSON object, counts its top-level keys, and injects
// the count into the configured field. The updated JSON is returned.
func (c *Counter) Apply(line []byte) ([]byte, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("fieldcount: invalid JSON: %w", err)
	}
	m[c.field] = json.RawMessage(fmt.Sprintf("%d", len(m)))
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("fieldcount: marshal: %w", err)
	}
	return out, nil
}
