package dropout

import (
	"fmt"
	"sync"
)

// Registry stores named Dropper instances.
type Registry struct {
	mu      sync.RWMutex
	droppers map[string]*Dropper
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{droppers: make(map[string]*Dropper)}
}

// Register adds a Dropper under name. Returns an error if the name is
// already taken.
func (r *Registry) Register(name string, d *Dropper) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.droppers[name]; exists {
		return fmt.Errorf("dropout: registry: %q already registered", name)
	}
	r.droppers[name] = d
	return nil
}

// Get retrieves a Dropper by name.
func (r *Registry) Get(name string) (*Dropper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.droppers[name]
	if !ok {
		return nil, fmt.Errorf("dropout: registry: %q not found", name)
	}
	return d, nil
}

// Names returns the sorted list of registered dropper names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.droppers))
	for n := range r.droppers {
		names = append(names, n)
	}
	return names
}
