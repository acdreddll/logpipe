// Package redact provides field-level redaction for structured log lines.
package redact

import (
	"encoding/json"
	"fmt"
)

// Redactor masks or removes sensitive fields from a JSON log line.
type Redactor struct {
	fields []string
	mask   string
}

// New creates a Redactor that replaces each named field value with mask.
// If mask is empty, the field is deleted from the object.
func New(fields []string, mask string) *Redactor {
	return &Redactor{fields: fields, mask: mask}
}

// Apply processes a raw JSON line and returns the redacted JSON.
func (r *Redactor) Apply(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("redact: invalid JSON: %w", err)
	}

	for _, f := range r.fields {
		if _, ok := obj[f]; !ok {
			continue
		}
		if r.mask == "" {
			delete(obj, f)
		} else {
			obj[f] = r.mask
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("redact: marshal: %w", err)
	}
	return string(out), nil
}
