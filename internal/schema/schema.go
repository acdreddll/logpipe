// Package schema provides JSON log schema validation against a set of
// required fields and expected types.
package schema

import (
	"encoding/json"
	"fmt"
)

// FieldType represents the expected JSON value type for a field.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeNumber  FieldType = "number"
	TypeBoolean FieldType = "boolean"
	TypeAny     FieldType = "any"
)

// FieldRule describes a single field constraint.
type FieldRule struct {
	Name     string
	Required bool
	Type     FieldType
}

// Validator checks log lines against a set of field rules.
type Validator struct {
	rules []FieldRule
}

// New creates a Validator from the provided rules.
// Returns an error if any rule has an empty name or unknown type.
func New(rules []FieldRule) (*Validator, error) {
	for _, r := range rules {
		if r.Name == "" {
			return nil, fmt.Errorf("schema: rule has empty field name")
		}
		switch r.Type {
		case TypeString, TypeNumber, TypeBoolean, TypeAny:
		default:
			return nil, fmt.Errorf("schema: unknown type %q for field %q", r.Type, r.Name)
		}
	}
	return &Validator{rules: rules}, nil
}

// Validate checks a raw JSON log line against the validator's rules.
// Returns a non-nil error describing the first violation found.
func (v *Validator) Validate(line []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return fmt.Errorf("schema: invalid JSON: %w", err)
	}

	for _, r := range v.rules {
		val, exists := obj[r.Name]
		if !exists {
			if r.Required {
				return fmt.Errorf("schema: required field %q is missing", r.Name)
			}
			continue
		}
		if r.Type == TypeAny {
			continue
		}
		if err := checkType(r.Name, r.Type, val); err != nil {
			return err
		}
	}
	return nil
}

func checkType(name string, expected FieldType, val any) error {
	switch expected {
	case TypeString:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("schema: field %q expected string, got %T", name, val)
		}
	case TypeNumber:
		if _, ok := val.(float64); !ok {
			return fmt.Errorf("schema: field %q expected number, got %T", name, val)
		}
	case TypeBoolean:
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("schema: field %q expected boolean, got %T", name, val)
		}
	}
	return nil
}
