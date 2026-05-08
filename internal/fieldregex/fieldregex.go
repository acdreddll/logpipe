package fieldregex

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Replacer replaces substrings in a JSON field value that match a regular
// expression with a fixed replacement string.
type Replacer struct {
	field       string
	pattern     *regexp.Regexp
	replacement string
}

// New creates a Replacer that applies the compiled pattern to the named field.
// An empty field name, empty pattern, or invalid regexp returns an error.
func New(field, pattern, replacement string) (*Replacer, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldregex: field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("fieldregex: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("fieldregex: invalid pattern: %w", err)
	}
	return &Replacer{field: field, pattern: re, replacement: replacement}, nil
}

// Apply parses line as JSON, replaces all matches of the pattern in the
// named field with the replacement string, and returns the modified JSON.
// If the field is absent the line is returned unchanged.
func (r *Replacer) Apply(line string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", fmt.Errorf("fieldregex: invalid JSON: %w", err)
	}

	v, ok := m[r.field]
	if !ok {
		return line, nil
	}

	s, ok := v.(string)
	if !ok {
		return line, nil
	}

	m[r.field] = r.pattern.ReplaceAllString(s, r.replacement)

	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("fieldregex: marshal: %w", err)
	}
	return string(out), nil
}
