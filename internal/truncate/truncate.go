package truncate

import (
	"encoding/json"
	"fmt"
)

// Truncator trims string field values that exceed a maximum byte length.
type Truncator struct {
	field  string
	maxLen int
	suffix string
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithSuffix sets the string appended to truncated values (default: "...").
func WithSuffix(s string) Option {
	return func(t *Truncator) { t.suffix = s }
}

// New creates a Truncator for the given field and maximum byte length.
// maxLen must be greater than zero.
func New(field string, maxLen int, opts ...Option) (*Truncator, error) {
	if field == "" {
		return nil, fmt.Errorf("truncate: field name must not be empty")
	}
	if maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be greater than zero")
	}
	t := &Truncator{
		field:  field,
		maxLen: maxLen,
		suffix: "...",
	}
	for _, o := range opts {
		o(t)
	}
	return t, nil
}

// Apply truncates the configured field in the JSON log line if its string
// value exceeds maxLen bytes. Non-string fields and missing fields are left
// unchanged. Returns the (possibly modified) JSON line.
func (t *Truncator) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("truncate: invalid JSON: %w", err)
	}

	v, ok := m[t.field]
	if !ok {
		return line, nil
	}

	s, ok := v.(string)
	if !ok {
		return line, nil
	}

	if len(s) > t.maxLen {
		truncated := s[:t.maxLen] + t.suffix
		m[t.field] = truncated
		out, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("truncate: marshal error: %w", err)
		}
		return out, nil
	}

	return line, nil
}
