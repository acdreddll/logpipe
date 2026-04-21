package masking

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Masker applies a regex-based masking pattern to a specific field in a JSON log line.
type Masker struct {
	field   string
	pattern *regexp.Regexp
	replace string
}

// New creates a Masker that replaces matches of pattern in field with replace.
// If replace is empty, it defaults to "***".
func New(field, pattern, replace string) (*Masker, error) {
	if field == "" {
		return nil, fmt.Errorf("masking: field must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("masking: invalid pattern %q: %w", pattern, err)
	}
	if replace == "" {
		replace = "***"
	}
	return &Masker{field: field, pattern: re, replace: replace}, nil
}

// Apply masks the target field in the JSON line and returns the modified line.
// If the field is absent the line is returned unchanged.
func (m *Masker) Apply(line []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("masking: invalid JSON: %w", err)
	}

	val, ok := obj[m.field]
	if !ok {
		return line, nil
	}

	str, ok := val.(string)
	if !ok {
		return line, nil
	}

	obj[m.field] = m.pattern.ReplaceAllString(str, m.replace)

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("masking: marshal error: %w", err)
	}
	return out, nil
}
