// Package typecast provides field type coercion for structured log lines.
// It converts string values in JSON log entries to target Go types such as
// int, float, or bool, leaving the rest of the document unchanged.
package typecast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Type represents a supported target type for coercion.
type Type string

const (
	TypeInt   Type = "int"
	TypeFloat Type = "float"
	TypeBool  Type = "bool"
	TypeString Type = "string"
)

// Caster coerces a named field in a JSON log line to the configured type.
type Caster struct {
	field  string
	target Type
}

// New creates a Caster that will coerce field to target.
// Returns an error if field is empty or target is unrecognised.
func New(field string, target Type) (*Caster, error) {
	if field == "" {
		return nil, fmt.Errorf("typecast: field name must not be empty")
	}
	switch target {
	case TypeInt, TypeFloat, TypeBool, TypeString:
	default:
		return nil, fmt.Errorf("typecast: unknown target type %q", target)
	}
	return &Caster{field: field, target: target}, nil
}

// Apply parses line as JSON, coerces the configured field, and returns the
// re-serialised line. If the field is absent the line is returned unchanged.
func (c *Caster) Apply(line []byte) ([]byte, error) {
	var doc map[string]interface{}
	if err := json.Unmarshal(line, &doc); err != nil {
		return nil, fmt.Errorf("typecast: invalid JSON: %w", err)
	}

	raw, ok := doc[c.field]
	if !ok {
		return line, nil
	}

	coerced, err := coerce(raw, c.target)
	if err != nil {
		return nil, fmt.Errorf("typecast: field %q: %w", c.field, err)
	}
	doc[c.field] = coerced

	out, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("typecast: marshal: %w", err)
	}
	return out, nil
}

func coerce(v interface{}, t Type) (interface{}, error) {
	s := fmt.Sprintf("%v", v)
	switch t {
	case TypeInt:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to int", s)
		}
		return n, nil
	case TypeFloat:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to float", s)
		}
		return f, nil
	case TypeBool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to bool", s)
		}
		return b, nil
	case TypeString:
		return s, nil
	}
	return nil, fmt.Errorf("unsupported type %q", t)
}
