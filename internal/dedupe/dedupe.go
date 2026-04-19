// Package dedupe provides log line deduplication based on a configurable field or full-line hashing.
package dedupe

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
)

// Deduper tracks seen log lines and drops duplicates within a fixed-size window.
type Deduper struct {
	fields []string
	mu     sync.Mutex
	seen   map[string]struct{}
	max    int
}

// New creates a Deduper that deduplicates based on the given fields.
// If fields is empty, the entire JSON line is used as the key.
// max controls the maximum number of unique keys retained before the cache is reset.
func New(fields []string, max int) (*Deduper, error) {
	if max <= 0 {
		return nil, fmt.Errorf("dedupe: max must be greater than zero")
	}
	return &Deduper{
		fields: fields,
		seen:   make(map[string]struct{}, max),
		max:    max,
	}, nil
}

// IsDuplicate returns true if the line has been seen before.
func (d *Deduper) IsDuplicate(line []byte) (bool, error) {
	key, err := d.keyFor(line)
	if err != nil {
		return false, err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.seen[key]; ok {
		return true, nil
	}
	if len(d.seen) >= d.max {
		d.seen = make(map[string]struct{}, d.max)
	}
	d.seen[key] = struct{}{}
	return false, nil
}

func (d *Deduper) keyFor(line []byte) (string, error) {
	if len(d.fields) == 0 {
		h := sha256.Sum256(line)
		return fmt.Sprintf("%x", h), nil
	}
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return "", fmt.Errorf("dedupe: invalid JSON: %w", err)
	}
	parts := make(map[string]any, len(d.fields))
	for _, f := range d.fields {
		parts[f] = obj[f]
	}
	b, err := json.Marshal(parts)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h), nil
}
