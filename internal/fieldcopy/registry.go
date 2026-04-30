package fieldcopy

import (
	"fmt"
	"sort"
)

// Registry holds named Copier instances.
type Registry struct {
	m map[string]*Copier
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{m: make(map[string]*Copier)}
}

// Register adds a Copier under name. Returns an error if the name is already taken.
func (r *Registry) Register(name string, c *Copier) error {
	if _, ok := r.m[name]; ok {
		return fmt.Errorf("fieldcopy: registry: %q already registered", name)
	}
	r.m[name] = c
	return nil
}

// Get retrieves a Copier by name.
func (r *Registry) Get(name string) (*Copier, error) {
	c, ok := r.m[name]
	if !ok {
		return nil, fmt.Errorf("fieldcopy: registry: %q not found", name)
	}
	return c, nil
}

// Names returns a sorted list of registered copier names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.m))
	for n := range r.m {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
