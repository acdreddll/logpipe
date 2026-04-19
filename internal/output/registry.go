package output

import (
	"fmt"
	"sync"
)

// Registry holds named output writers.
type Registry struct {
	mu      sync.RWMutex
	writers map[string]*Writer
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{writers: make(map[string]*Writer)}
}

// Register adds a Writer to the registry. Returns error if name already exists.
func (r *Registry) Register(w *Writer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.writers[w.Name]; exists {
		return fmt.Errorf("output: writer %q already registered", w.Name)
	}
	r.writers[w.Name] = w
	return nil
}

// Get retrieves a Writer by name.
func (r *Registry) Get(name string) (*Writer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.writers[name]
	return w, ok
}

// Names returns all registered writer names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.writers))
	for n := range r.writers {
		names = append(names, n)
	}
	return names
}

// Remove unregisters a writer by name.
func (r *Registry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.writers, name)
}
