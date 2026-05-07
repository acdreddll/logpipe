package fieldround

import (
	"fmt"
)

// Registry holds named Rounder instances.
type Registry struct {
	items map[string]*Rounder
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Rounder)}
}

// Register adds a Rounder under name.
func (reg *Registry) Register(name string, r *Rounder) error {
	if _, exists := reg.items[name]; exists {
		return fmt.Errorf("fieldround: registry: %q already registered", name)
	}
	reg.items[name] = r
	return nil
}

// Get retrieves a Rounder by name.
func (reg *Registry) Get(name string) (*Rounder, error) {
	r, ok := reg.items[name]
	if !ok {
		return nil, fmt.Errorf("fieldround: registry: %q not found", name)
	}
	return r, nil
}

// Names returns the sorted list of registered names.
func (reg *Registry) Names() []string {
	out := make([]string, 0, len(reg.items))
	for k := range reg.items {
		out = append(out, k)
	}
	return out
}
