// Package masking provides regex-based field masking for structured JSON log lines.
//
// A Masker targets a single named field and replaces all substrings matching
// a compiled regular expression with a configurable replacement string.
// This is useful for scrubbing PII such as email addresses, tokens, or
// credit-card numbers from log streams before they are routed to outputs.
//
// # Masker
//
// Each Masker is immutable after construction and safe for concurrent use.
// The Apply method accepts a raw JSON-encoded log line as a byte slice and
// returns a new byte slice with the field value masked. If the target field
// is absent from the log line, the input is returned unchanged.
//
// # Example usage
//
//	m, err := masking.New("email", `[^@]+@[^@]+`, "[REDACTED]")
//	if err != nil { ... }
//	masked, err := m.Apply(line)
//
// Multiple Maskers can be composed to redact several fields in a single pass:
//
//	maskers := []*masking.Masker{
//		emailMasker,
//		tokenMasker,
//	}
//	for _, m := range maskers {
//		line, err = m.Apply(line)
//		if err != nil { ... }
//	}
//
// # ApplyAll helper
//
// As a convenience, ApplyAll applies a slice of Maskers to a log line in
// order, short-circuiting and returning the first error encountered:
//
//	line, err := masking.ApplyAll(line, emailMasker, tokenMasker)
package masking
