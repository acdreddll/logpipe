package ratelimit

import (
	"fmt"
	"sync"
)

// Registry holds named rate limiters, keyed by route or output name.
type Registry struct {
	mu       sync.RWMutex
	limiters map[string]*Limiter
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{limiters: make(map[string]*Limiter)}
}

// Register adds a Limiter with the given name and rate.
// Returns an error if the name is already registered or rate is invalid.
func (r *Registry) Register(name string, ratePerSec float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.limiters[name]; exists {
		return fmt.Errorf("ratelimit: limiter %q already registered", name)
	}
	l, err := New(ratePerSec)
	if err != nil {
		return err
	}
	r.limiters[name] = l
	return nil
}

// Get returns the Limiter registered under name, or an error if not found.
func (r *Registry) Get(name string) (*Limiter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, ok := r.limiters[name]
	if !ok {
		return nil, fmt.Errorf("ratelimit: no limiter registered for %q", name)
	}
	return l, nil
}

// Names returns all registered limiter names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.limiters))
	for n := range r.limiters {
		names = append(names, n)
	}
	return names
}
