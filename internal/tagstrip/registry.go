package tagstrip

import (
	"fmt"
)

// Registry holds named Stripper instances.
type Registry struct {
	strippers map[string]*Stripper
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{strippers: make(map[string]*Stripper)}
}

// Register adds a named Stripper built from the given fields.
// Returns an error if the name is already registered or construction fails.
func (r *Registry) Register(name string, fields []string) error {
	if _, exists := r.strippers[name]; exists {
		return fmt.Errorf("tagstrip: registry: %q already registered", name)
	}
	s, err := New(fields)
	if err != nil {
		return fmt.Errorf("tagstrip: registry: %w", err)
	}
	r.strippers[name] = s
	return nil
}

// Get retrieves a Stripper by name.
func (r *Registry) Get(name string) (*Stripper, error) {
	s, ok := r.strippers[name]
	if !ok {
		return nil, fmt.Errorf("tagstrip: registry: %q not found", name)
	}
	return s, nil
}

// Names returns all registered Stripper names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.strippers))
	for n := range r.strippers {
		names = append(names, n)
	}
	return names
}
