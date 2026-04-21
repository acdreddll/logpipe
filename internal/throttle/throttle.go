package throttle

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Throttler suppresses repeated identical log lines within a cooldown window.
// Once a line is seen, subsequent identical lines are dropped until the
// cooldown duration elapses.
type Throttler struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
	now      func() time.Time
}

// New creates a Throttler with the given cooldown duration.
// cooldown must be greater than zero.
func New(cooldown time.Duration) (*Throttler, error) {
	if cooldown <= 0 {
		return nil, fmt.Errorf("throttle: cooldown must be greater than zero")
	}
	return &Throttler{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
		now:      time.Now,
	}, nil
}

// Allow returns true if the log line should be forwarded, false if it should
// be suppressed. A line is suppressed when an identical line was already
// forwarded within the cooldown window.
func (t *Throttler) Allow(line []byte) (bool, error) {
	key, err := normalise(line)
	if err != nil {
		return false, fmt.Errorf("throttle: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.seen[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false, nil
		}
	}
	t.seen[key] = now
	return true, nil
}

// Purge removes expired entries from the internal tracking map.
// Call periodically to prevent unbounded memory growth.
func (t *Throttler) Purge() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	for k, ts := range t.seen {
		if now.Sub(ts) >= t.cooldown {
			delete(t.seen, k)
		}
	}
}

// normalise returns a canonical string key for a JSON log line by
// round-tripping through map[string]any so that key ordering is stable.
func normalise(line []byte) (string, error) {
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	out, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
