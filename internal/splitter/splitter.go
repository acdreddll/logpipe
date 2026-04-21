package splitter

import (
	"encoding/json"
	"fmt"
)

// Splitter fans out a single log line to multiple named outputs based on a
// field value. Lines that do not contain the routing field are sent to the
// default bucket (empty string key).
type Splitter struct {
	field    string
	routes   map[string][]string // field value -> bucket names
	default_ string
}

// Option configures a Splitter.
type Option func(*Splitter)

// WithDefault sets the bucket name used when no route matches.
func WithDefault(name string) Option {
	return func(s *Splitter) { s.default_ = name }
}

// New creates a Splitter that inspects the given field and maps its values to
// bucket names via the routes map.
func New(field string, routes map[string][]string, opts ...Option) (*Splitter, error) {
	if field == "" {
		return nil, fmt.Errorf("splitter: field must not be empty")
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("splitter: routes must not be empty")
	}
	s := &Splitter{
		field:  field,
		routes: routes,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// Split returns the list of bucket names the line should be forwarded to.
// An empty slice means the line should be dropped.
func (s *Splitter) Split(line []byte) ([]string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("splitter: invalid JSON: %w", err)
	}

	v, ok := obj[s.field]
	if !ok {
		if s.default_ != "" {
			return []string{s.default_}, nil
		}
		return nil, nil
	}

	key := fmt.Sprintf("%v", v)
	buckets, matched := s.routes[key]
	if !matched {
		if s.default_ != "" {
			return []string{s.default_}, nil
		}
		return nil, nil
	}
	return buckets, nil
}

// Field returns the routing field name.
func (s *Splitter) Field() string { return s.field }
