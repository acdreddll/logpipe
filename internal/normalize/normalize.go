// Package normalize provides field-level key normalization for structured log lines.
// It renames or lowercases JSON keys according to a configured mapping,
// making downstream filtering and routing more predictable.
package normalize

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Normalizer applies key normalization to JSON log lines.
type Normalizer struct {
	// mapping holds explicit old-key → new-key renames.
	mapping map[string]string
	// lowercase, when true, lowercases every key not covered by mapping.
	lowercase bool
}

// Option configures a Normalizer.
type Option func(*Normalizer)

// WithMapping adds explicit key rename pairs.
func WithMapping(m map[string]string) Option {
	return func(n *Normalizer) {
		for k, v := range m {
			n.mapping[k] = v
		}
	}
}

// WithLowercase enables automatic lowercasing of all keys.
func WithLowercase() Option {
	return func(n *Normalizer) {
		n.lowercase = true
	}
}

// New creates a Normalizer with the given options.
// Returns an error if no options are provided.
func New(opts ...Option) (*Normalizer, error) {
	if len(opts) == 0 {
		return nil, fmt.Errorf("normalize: at least one option required")
	}
	n := &Normalizer{mapping: make(map[string]string)}
	for _, o := range opts {
		o(n)
	}
	return n, nil
}

// Apply normalizes keys in the JSON line and returns the updated line.
func (n *Normalizer) Apply(line string) (string, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("normalize: invalid JSON: %w", err)
	}

	out := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		newKey := n.resolveKey(k)
		out[newKey] = v
	}

	b, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("normalize: marshal error: %w", err)
	}
	return string(b), nil
}

func (n *Normalizer) resolveKey(k string) string {
	if renamed, ok := n.mapping[k]; ok {
		return renamed
	}
	if n.lowercase {
		return strings.ToLower(k)
	}
	return k
}
