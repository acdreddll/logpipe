// Package config loads and validates logpipe pipeline configuration from YAML.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RouteConfig defines a single named route with an optional filter and output.
type RouteConfig struct {
	Name      string            `yaml:"name"`
	Filter    map[string]string `yaml:"filter,omitempty"`
	Transform []map[string]string `yaml:"transform,omitempty"`
	Output    OutputConfig      `yaml:"output"`
}

// OutputConfig specifies the output destination for a route.
type OutputConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path,omitempty"`
}

// Config is the top-level pipeline configuration.
type Config struct {
	Routes []RouteConfig `yaml:"routes"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Routes) == 0 {
		return fmt.Errorf("config: at least one route is required")
	}
	for i, r := range c.Routes {
		if r.Name == "" {
			return fmt.Errorf("config: route[%d]: name is required", i)
		}
		if r.Output.Type == "" {
			return fmt.Errorf("config: route[%d] %q: output.type is required", i, r.Name)
		}
	}
	return nil
}
