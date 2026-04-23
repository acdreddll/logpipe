// Package retag provides a processor that renames or adds a tag field
// in a JSON log line based on a static mapping of old-value → new-value.
package retag

import (
	"encoding/json"
	"fmt"
)

// Retagger rewrites the value of a single field according to a lookup table.
type Retagger struct {
	field   string
	mapping map[string]string
	defaultTag string
}

// Option is a functional option for Retagger.
type Option func(*Retagger)

// WithDefault sets a fallback tag value when the field value is not in the mapping.
func WithDefault(tag string) Option {
	return func(r *Retagger) { r.defaultTag = tag }
}

// New creates a Retagger that rewrites r.field using mapping.
// field must be non-empty and mapping must contain at least one entry.
func New(field string, mapping map[string]string, opts ...Option) (*Retagger, error) {
	if field == "" {
		return nil, fmt.Errorf("retag: field must not be empty")
	}
	if len(mapping) == 0 {
		return nil, fmt.Errorf("retag: mapping must contain at least one entry")
	}
	r := &Retagger{field: field, mapping: mapping}
	for _, o := range opts {
		o(r)
	}
	return r, nil
}

// Apply rewrites the target field in line according to the mapping.
// Lines where the field is absent are returned unchanged.
// If the field value is not in the mapping and no default is set, the line is returned unchanged.
func (r *Retagger) Apply(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("retag: invalid JSON: %w", err)
	}

	raw, ok := obj[r.field]
	if !ok {
		return line, nil
	}

	current, _ := raw.(string)
	if next, found := r.mapping[current]; found {
		obj[r.field] = next
	} else if r.defaultTag != "" {
		obj[r.field] = r.defaultTag
	} else {
		return line, nil
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("retag: marshal error: %w", err)
	}
	return string(out), nil
}
