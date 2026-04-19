package alert

import (
	"fmt"
	"sync"
)

// Registry holds named alerts.
type Registry struct {
	mu     sync.RWMutex
	alerts map[string]*Alert
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{alerts: make(map[string]*Alert)}
}

// Register adds an alert to the registry.
func (r *Registry) Register(a *Alert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.alerts[a.Name]; exists {
		return fmt.Errorf("alert %q already registered", a.Name)
	}
	r.alerts[a.Name] = a
	return nil
}

// Get retrieves an alert by name.
func (r *Registry) Get(name string) (*Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.alerts[name]
	if !ok {
		return nil, fmt.Errorf("alert %q not found", name)
	}
	return a, nil
}

// EvaluateAll runs all registered alerts against the given line.
func (r *Registry) EvaluateAll(line string) []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var errs []error
	for _, a := range r.alerts {
		if err := a.Evaluate(line); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Names returns all registered alert names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.alerts))
	for n := range r.alerts {
		names = append(names, n)
	}
	return names
}
