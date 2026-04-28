// Package debounce suppresses repeated log events within a configurable
// quiet window, emitting only the first occurrence until the window expires.
package debounce

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Debouncer holds state for debouncing log lines by a key field.
type Debouncer struct {
	field   string
	window  time.Duration
	mu      sync.Mutex
	seen    map[string]time.Time
}

// New creates a Debouncer that suppresses repeated values of field within
// the given quiet window. window must be positive.
func New(field string, window time.Duration) (*Debouncer, error) {
	if field == "" {
		return nil, errors.New("debounce: field must not be empty")
	}
	if window <= 0 {
		return nil, errors.New("debounce: window must be positive")
	}
	return &Debouncer{
		field:  field,
		window: window,
		seen:   make(map[string]time.Time),
	}, nil
}

// Allow returns true if the log line should be forwarded. The first
// occurrence of a key value is always allowed; subsequent occurrences
// within the quiet window are suppressed. Lines with missing or
// non-string key values are always allowed.
func (d *Debouncer) Allow(line []byte) (bool, error) {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return false, err
	}

	v, ok := obj[d.field]
	if !ok {
		return true, nil
	}
	key, ok := v.(string)
	if !ok {
		return true, nil
	}

	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()

	if last, exists := d.seen[key]; exists && now.Sub(last) < d.window {
		return false, nil
	}
	d.seen[key] = now
	return true, nil
}

// Purge removes expired entries from the internal state map.
func (d *Debouncer) Purge() {
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}
