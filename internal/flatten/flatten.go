package flatten

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Flattener flattens nested JSON objects into a single-level map
// using a configurable separator between key segments.
type Flattener struct {
	separator string
	prefix    string
}

// Option configures a Flattener.
type Option func(*Flattener)

// WithSeparator sets the key separator (default: ".").
func WithSeparator(sep string) Option {
	return func(f *Flattener) {
		if sep != "" {
			f.separator = sep
		}
	}
}

// WithPrefix sets an optional prefix prepended to every top-level key.
func WithPrefix(prefix string) Option {
	return func(f *Flattener) {
		f.prefix = prefix
	}
}

// New creates a Flattener with the supplied options.
func New(opts ...Option) *Flattener {
	f := &Flattener{separator: "."}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Apply receives a JSON line, flattens any nested objects, and returns
// the re-encoded JSON line. Arrays are left intact.
func (f *Flattener) Apply(line string) (string, error) {
	var root map[string]any
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return "", fmt.Errorf("flatten: invalid JSON: %w", err)
	}

	flat := make(map[string]any)
	f.flatten(root, f.prefix, flat)

	out, err := json.Marshal(flat)
	if err != nil {
		return "", fmt.Errorf("flatten: marshal: %w", err)
	}
	return string(out), nil
}

func (f *Flattener) flatten(obj map[string]any, prefix string, out map[string]any) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, f.separator)
		}
		switch child := v.(type) {
		case map[string]any:
			f.flatten(child, key, out)
		default:
			out[key] = v
		}
	}
}
