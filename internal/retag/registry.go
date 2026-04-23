package retag

import (
	"fmt"
	"sync"
)

// Registry holds named Retagger instances.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]*Retagger
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]*Retagger)}
}

// Register adds a Retagger under name. Returns an error if the name is already taken.
func (reg *Registry) Register(name string, r *Retagger) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if _, ok := reg.entries[name]; ok {
		return fmt.Errorf("retag: registry: %q already registered", name)
	}
	reg.entries[name] = r
	return nil
}

// Get retrieves a Retagger by name.
func (reg *Registry) Get(name string) (*Retagger, error) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	r, ok := reg.entries[name]
	if !ok {
		return nil, fmt.Errorf("retag: registry: %q not found", name)
	}
	return r, nil
}

// Names returns the sorted list of registered names.
func (reg *Registry) Names() []string {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	names := make([]string, 0, len(reg.entries))
	for n := range reg.entries {
		names = append(names, n)
	}
	return names
}
