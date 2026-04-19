// Package config loads and validates logpipe YAML configuration.
package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Route describes a single named routing rule.
type Route struct {
	Name      string            `yaml:"name"`
	Filter    map[string]string `yaml:"filter"`
	Transform []map[string]string `yaml:"transform"`
	Output    string            `yaml:"output"`
	Sample    *SampleConfig     `yaml:"sample"`
	Redact    []RedactConfig    `yaml:"redact"`
	Buffer    *BufferConfig     `yaml:"buffer"`
}

// BufferConfig holds optional buffering settings for a route.
type BufferConfig struct {
	Size     int    `yaml:"size"`
	Interval string `yaml:"interval"`
}

// SampleConfig defines sampling behaviour for a route.
type SampleConfig struct {
	Rate float64 `yaml:"rate"`
	Mode string  `yaml:"mode"` // "deterministic" | "random"
}

// RedactConfig describes a single field redaction rule.
type RedactConfig struct {
	Field  string `yaml:"field"`
	Action string `yaml:"action"` // "mask" | "delete"
}

// OutputConfig describes an output destination.
type OutputConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

// Config is the top-level logpipe configuration.
type Config struct {
	Outputs []OutputConfig `yaml:"outputs"`
	Routes  []Route        `yaml:"routes"`
}

// Load reads and validates a YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	if len(cfg.Routes) == 0 {
		return errors.New("config: at least one route is required")
	}
	for _, o := range cfg.Outputs {
		if o.Type == "" {
			return fmt.Errorf("config: output %q missing type", o.Name)
		}
	}
	return nil
}
