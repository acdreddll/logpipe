// Package fieldmatch provides a processor that checks whether a log field
// value matches a regular expression and writes the boolean result into a
// configurable destination field.
package fieldmatch

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Matcher checks a source field against a compiled regular expression and
// injects the boolean result into a destination field.
type Matcher struct {
	src     string
	dest    string
	pattern *regexp.Regexp
}

// New returns a Matcher that reads src, tests it against pattern, and writes
// the boolean result into dest.
func New(src, dest, pattern string) (*Matcher, error) {
	if src == "" {
		return nil, fmt.Errorf("fieldmatch: src field must not be empty")
	}
	if dest == "" {
		return nil, fmt.Errorf("fieldmatch: dest field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("fieldmatch: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("fieldmatch: invalid pattern %q: %w", pattern, err)
	}
	return &Matcher{src: src, dest: dest, pattern: re}, nil
}

// Apply evaluates the matcher against line (a JSON object) and returns a new
// JSON line with the dest field set to true or false.
func (m *Matcher) Apply(line string) (string, error) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("fieldmatch: invalid JSON: %w", err)
	}

	matched := false
	if v, ok := obj[m.src]; ok {
		if s, ok := v.(string); ok {
			matched = m.pattern.MatchString(s)
		}
	}
	obj[m.dest] = matched

	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("fieldmatch: marshal: %w", err)
	}
	return string(out), nil
}
