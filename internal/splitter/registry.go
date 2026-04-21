package splitter

import (
	"fmt"
	"sync"
)

// Registry holds named Splitter instances.
type Registry struct {
	mu       sync.RWMutex
	entries  map[string]*Splitter
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]*Splitter)}
}

// Register adds a Splitter under the given name.
func (r *Registry) Register(name string, s *Splitter) error {
	if name == "" {
		return fmt.Errorf("splitter registry: name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("splitter registry: %q already registered", name)
	}
	r.entries[name] = s
	return nil
}

// Get retrieves a Splitter by name.
func (r *Registry) Get(name string) (*Splitter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.entries[name]
	if !ok {
		return nil, fmt.Errorf("splitter registry: %q not found", name)
	}
	return s, nil
}

// Names returns the sorted list of registered splitter names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.entries))
	for k := range r.entries {
		names = append(names, k)
	}
	return names
}
