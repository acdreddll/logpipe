package enrichment

import (
	"encoding/json"
	"fmt"
)

// Enricher adds static or derived fields to a log line.
type Enricher struct {
	fields map[string]string
}

// New creates an Enricher that will inject the given static fields into every log line.
func New(fields map[string]string) (*Enricher, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("enrichment: at least one field is required")
	}
	return &Enricher{fields: fields}, nil
}

// Apply merges the configured fields into the JSON log line.
// Existing keys are NOT overwritten.
func (e *Enricher) Apply(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("enrichment: invalid JSON: %w", err)
	}
	for k, v := range e.fields {
		if _, exists := obj[k]; !exists {
			obj[k] = v
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("enrichment: marshal error: %w", err)
	}
	return string(out), nil
}

// Fields returns a copy of the static fields configured on the Enricher.
func (e *Enricher) Fields() map[string]string {
	copy := make(map[string]string, len(e.fields))
	for k, v := range e.fields {
		copy[k] = v
	}
	return copy
}
