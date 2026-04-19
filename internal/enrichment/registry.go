package enrichment

import (
	"fmt"
	"sync"
)

// Registry holds named Enricher instances.
type Registry struct {
	mu       sync.RWMutex
	enrichers map[string]*Enricher
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{enrichers: make(map[string]*Enricher)}
}

// Register adds a named Enricher. Returns an error if the name is already taken.
func (r *Registry) Register(name string, e *Enricher) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.enrichers[name]; exists {
		return fmt.Errorf("enrichment: enricher %q already registered", name)
	}
	r.enrichers[name] = e
	return nil
}

// Get retrieves a named Enricher.
func (r *Registry) Get(name string) (*Enricher, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.enrichers[name]
	if !ok {
		return nil, fmt.Errorf("enrichment: enricher %q not found", name)
	}
	return e, nil
}

// Names returns all registered enricher names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.enrichers))
	for n := range r.enrichers {
		names = append(names, n)
	}
	return names
}
