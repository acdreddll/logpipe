package flatten

import (
	"fmt"
	"sort"
)

// Registry holds named Flattener instances.
type Registry struct {
	items map[string]*Flattener
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Flattener)}
}

// Register adds a Flattener under name. Returns an error if the name is
// already taken.
func (r *Registry) Register(name string, f *Flattener) error {
	if _, exists := r.items[name]; exists {
		return fmt.Errorf("flatten: registry: %q already registered", name)
	}
	r.items[name] = f
	return nil
}

// Get returns the Flattener registered under name.
func (r *Registry) Get(name string) (*Flattener, error) {
	f, ok := r.items[name]
	if !ok {
		return nil, fmt.Errorf("flatten: registry: %q not found", name)
	}
	return f, nil
}

// Names returns all registered names in sorted order.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.items))
	for n := range r.items {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
