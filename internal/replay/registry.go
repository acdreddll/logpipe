package replay

import (
	"fmt"
	"sync"
)

// Registry stores named Replayer instances.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]*Replayer
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]*Replayer)}
}

// Register adds a Replayer under name. Returns an error if the name is already
// taken.
func (reg *Registry) Register(name string, r *Replayer) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if _, ok := reg.entries[name]; ok {
		return fmt.Errorf("replay: %q already registered", name)
	}
	reg.entries[name] = r
	return nil
}

// Get retrieves a Replayer by name.
func (reg *Registry) Get(name string) (*Replayer, error) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	r, ok := reg.entries[name]
	if !ok {
		return nil, fmt.Errorf("replay: %q not found", name)
	}
	return r, nil
}

// Names returns all registered names in insertion order (map order).
func (reg *Registry) Names() []string {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	names := make([]string, 0, len(reg.entries))
	for n := range reg.entries {
		names = append(names, n)
	}
	return names
}
