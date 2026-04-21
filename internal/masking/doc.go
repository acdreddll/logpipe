// Package masking provides regex-based field masking for structured JSON log lines.
//
// A Masker targets a single named field and replaces all substrings matching
// a compiled regular expression with a configurable replacement string.
// This is useful for scrubbing PII such as email addresses, tokens, or
// credit-card numbers from log streams before they are routed to outputs.
//
// Example usage:
//
//	m, err := masking.New("email", `[^@]+@[^@]+`, "[REDACTED]")
//	if err != nil { ... }
//	masked, err := m.Apply(line)
package masking
