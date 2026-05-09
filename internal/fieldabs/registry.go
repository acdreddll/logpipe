package fieldabs

import "fmt"

// Registry holds named Processor instances.
type Registry struct {
	processors map[string]*Processor
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{processors: make(map[string]*Processor)}
}

// Register adds a Processor under the given name.
// Returns an error on duplicate names or construction failure.
func (r *Registry) Register(name, field string) error {
	if _, exists := r.processors[name]; exists {
		return fmt.Errorf("fieldabs: processor %q already registered", name)
	}
	p, err := New(field)
	if err != nil {
		return err
	}
	r.processors[name] = p
	return nil
}

// Get returns the Processor registered under name, or an error if not found.
func (r *Registry) Get(name string) (*Processor, error) {
	p, ok := r.processors[name]
	if !ok {
		return nil, fmt.Errorf("fieldabs: processor %q not found", name)
	}
	return p, nil
}

// Names returns all registered processor names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.processors))
	for n := range r.processors {
		names = append(names, n)
	}
	return names
}
