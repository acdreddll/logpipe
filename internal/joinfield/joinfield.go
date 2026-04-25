// Package joinfield provides a processor that concatenates multiple source
// fields into a single destination field using a configurable separator.
package joinfield

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Joiner concatenates values from several source fields into one destination
// field. Missing or non-string source fields are silently skipped.
type Joiner struct {
	sources []string
	dest    string
	sep     string
}

// WithSeparator overrides the default separator (a single space).
func WithSeparator(sep string) func(*Joiner) {
	return func(j *Joiner) { j.sep = sep }
}

// New creates a Joiner that reads from sources and writes to dest.
// At least one source field and a non-empty dest are required.
func New(dest string, sources []string, opts ...func(*Joiner)) (*Joiner, error) {
	if dest == "" {
		return nil, errors.New("joinfield: dest field must not be empty")
	}
	if len(sources) == 0 {
		return nil, errors.New("joinfield: at least one source field is required")
	}
	for _, s := range sources {
		if s == "" {
			return nil, errors.New("joinfield: source field names must not be empty")
		}
	}
	j := &Joiner{dest: dest, sources: sources, sep: " "}
	for _, o := range opts {
		o(j)
	}
	return j, nil
}

// Apply joins the source fields in line and returns the modified JSON.
func (j *Joiner) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("joinfield: invalid JSON: %w", err)
	}

	var parts []string
	for _, src := range j.sources {
		v, ok := m[src]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if s != "" {
			parts = append(parts, s)
		}
	}

	m[j.dest] = strings.Join(parts, j.sep)

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("joinfield: marshal error: %w", err)
	}
	return out, nil
}
