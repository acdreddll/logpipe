package transform

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds a named list of transform rules, suitable for loading
// from a JSON config file.
type Config struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

// Registry maps transformer names to Transformer instances.
type Registry struct {
	transformers map[string]*Transformer
}

// NewRegistry builds a Registry from a slice of Config entries.
func NewRegistry(configs []Config) (*Registry, error) {
	reg := &Registry{transformers: make(map[string]*Transformer, len(configs))}
	for _, cfg := range configs {
		if cfg.Name == "" {
			return nil, fmt.Errorf("transform config missing name")
		}
		tr, err := New(cfg.Rules)
		if err != nil {
			return nil, fmt.Errorf("transform %q: %w", cfg.Name, err)
		}
		reg.transformers[cfg.Name] = tr
	}
	return reg, nil
}

// Get returns the named Transformer or an error if not found.
func (r *Registry) Get(name string) (*Transformer, error) {
	tr, ok := r.transformers[name]
	if !ok {
		return nil, fmt.Errorf("transform %q not found", name)
	}
	return tr, nil
}

// LoadFile reads a JSON file containing an array of Config and returns a Registry.
func LoadFile(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading transform config: %w", err)
	}
	var configs []Config
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("parsing transform config: %w", err)
	}
	return NewRegistry(configs)
}
