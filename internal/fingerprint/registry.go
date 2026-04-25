package fingerprint

import (
	"fmt"
	"sort"
)

// Registry stores named Fingerprinter instances.
type Registry struct {
	entries map[string]*Fingerprinter
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]*Fingerprinter)}
}

// Register adds a Fingerprinter under name. Returns an error if the name is
// already registered.
func (r *Registry) Register(name string, fp *Fingerprinter) error {
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("fingerprint: %q already registered", name)
	}
	r.entries[name] = fp
	return nil
}

// Get returns the Fingerprinter registered under name.
func (r *Registry) Get(name string) (*Fingerprinter, error) {
	fp, ok := r.entries[name]
	if !ok {
		return nil, fmt.Errorf("fingerprint: %q not found", name)
	}
	return fp, nil
}

// Names returns a sorted list of registered names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.entries))
	for n := range r.entries {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
