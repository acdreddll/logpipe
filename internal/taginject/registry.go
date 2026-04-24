package taginject

import (
	"fmt"
	"sort"
)

// Registry stores named Injector instances.
type Registry struct {
	m map[string]*Injector
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{m: make(map[string]*Injector)}
}

// Register adds an Injector under name.
// Returns an error if the name is already registered.
func (r *Registry) Register(name string, inj *Injector) error {
	if _, ok := r.m[name]; ok {
		return fmt.Errorf("taginject: %q already registered", name)
	}
	r.m[name] = inj
	return nil
}

// Get retrieves a registered Injector by name.
func (r *Registry) Get(name string) (*Injector, error) {
	inj, ok := r.m[name]
	if !ok {
		return nil, fmt.Errorf("taginject: %q not found", name)
	}
	return inj, nil
}

// Names returns sorted names of all registered injectors.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.m))
	for k := range r.m {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
