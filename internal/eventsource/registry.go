package eventsource

import (
	"fmt"
	"sort"
)

// Registry holds named Source instances.
type Registry struct {
	sources map[string]*Source
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{sources: make(map[string]*Source)}
}

// Register adds a Source under the given key.
// Returns an error if the key is already registered.
func (r *Registry) Register(key string, s *Source) error {
	if _, ok := r.sources[key]; ok {
		return fmt.Errorf("eventsource: registry: %q already registered", key)
	}
	r.sources[key] = s
	return nil
}

// Get retrieves a Source by key.
func (r *Registry) Get(key string) (*Source, error) {
	s, ok := r.sources[key]
	if !ok {
		return nil, fmt.Errorf("eventsource: registry: %q not found", key)
	}
	return s, nil
}

// Names returns a sorted list of registered keys.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.sources))
	for k := range r.sources {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
