package debounce

import (
	"errors"
	"fmt"
	"time"
)

// Registry stores named Debouncer instances.
type Registry struct {
	items map[string]*Debouncer
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Debouncer)}
}

// Register adds a new Debouncer under name. Returns an error if the name
// is already registered or if the parameters are invalid.
func (r *Registry) Register(name, field string, window time.Duration) error {
	if name == "" {
		return errors.New("debounce registry: name must not be empty")
	}
	if _, exists := r.items[name]; exists {
		return fmt.Errorf("debounce registry: %q already registered", name)
	}
	d, err := New(field, window)
	if err != nil {
		return fmt.Errorf("debounce registry: %w", err)
	}
	r.items[name] = d
	return nil
}

// Get retrieves a Debouncer by name.
func (r *Registry) Get(name string) (*Debouncer, error) {
	d, ok := r.items[name]
	if !ok {
		return nil, fmt.Errorf("debounce registry: %q not found", name)
	}
	return d, nil
}

// Names returns the registered Debouncer names in insertion order.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.items))
	for n := range r.items {
		names = append(names, n)
	}
	return names
}
