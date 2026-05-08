package fieldregex

import "fmt"

// Registry holds named Replacer instances.
type Registry struct {
	items map[string]*Replacer
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Replacer)}
}

// Register adds a Replacer under name. Duplicate names return an error.
func (reg *Registry) Register(name string, r *Replacer) error {
	if _, exists := reg.items[name]; exists {
		return fmt.Errorf("fieldregex: registry: %q already registered", name)
	}
	reg.items[name] = r
	return nil
}

// Get returns the Replacer registered under name, or an error if absent.
func (reg *Registry) Get(name string) (*Replacer, error) {
	r, ok := reg.items[name]
	if !ok {
		return nil, fmt.Errorf("fieldregex: registry: %q not found", name)
	}
	return r, nil
}

// Names returns all registered names in insertion-independent order.
func (reg *Registry) Names() []string {
	out := make([]string, 0, len(reg.items))
	for k := range reg.items {
		out = append(out, k)
	}
	return out
}
