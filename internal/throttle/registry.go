package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Registry holds named Throttler instances.
type Registry struct {
	mu    sync.RWMutex
	items map[string]*Throttler
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Throttler)}
}

// Register creates and stores a Throttler under name with the given cooldown.
func (r *Registry) Register(name string, cooldown time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[name]; exists {
		return fmt.Errorf("throttle: %q already registered", name)
	}
	th, err := New(cooldown)
	if err != nil {
		return err
	}
	r.items[name] = th
	return nil
}

// Get retrieves a Throttler by name.
func (r *Registry) Get(name string) (*Throttler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	th, ok := r.items[name]
	if !ok {
		return nil, fmt.Errorf("throttle: %q not found", name)
	}
	return th, nil
}

// Names returns all registered throttler names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.items))
	for n := range r.items {
		names = append(names, n)
	}
	return names
}
