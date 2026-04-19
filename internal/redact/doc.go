// Package redact implements field-level redaction for structured JSON log lines.
//
// A Redactor is configured with a list of field names and an optional mask string.
// When applied to a JSON log line, each named field is either replaced with the
// mask value or deleted entirely if the mask is empty.
//
// Example usage:
//
//	r := redact.New([]string{"password", "token"}, "***")
//	clean, err := r.Apply(rawLine)
package redact
