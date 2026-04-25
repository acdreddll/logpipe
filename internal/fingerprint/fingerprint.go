// Package fingerprint computes a stable hash fingerprint for a log event,
// optionally restricted to a subset of fields. The fingerprint can be used
// for deduplication, correlation, or routing decisions.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// Fingerprinter computes a deterministic hex fingerprint for a JSON log line.
type Fingerprinter struct {
	fields []string // empty means hash the whole event
}

// Option is a functional option for Fingerprinter.
type Option func(*Fingerprinter)

// WithFields restricts the fingerprint to the named top-level fields.
func WithFields(fields ...string) Option {
	return func(f *Fingerprinter) {
		f.fields = append(f.fields, fields...)
	}
}

// New creates a Fingerprinter. At least one option may be supplied.
func New(opts ...Option) (*Fingerprinter, error) {
	fp := &Fingerprinter{}
	for _, o := range opts {
		o(fp)
	}
	return fp, nil
}

// Compute returns the SHA-256 fingerprint of the log line.
// If fields were specified only those fields contribute to the hash.
// Returns an error if line is not valid JSON.
func (fp *Fingerprinter) Compute(line []byte) (string, error) {
	var event map[string]json.RawMessage
	if err := json.Unmarshal(line, &event); err != nil {
		return "", fmt.Errorf("fingerprint: invalid JSON: %w", err)
	}

	var src map[string]json.RawMessage
	if len(fp.fields) == 0 {
		src = event
	} else {
		src = make(map[string]json.RawMessage, len(fp.fields))
		for _, f := range fp.fields {
			if v, ok := event[f]; ok {
				src[f] = v
			}
		}
	}

	// Sort keys for determinism.
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s;", k, src[k])
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
