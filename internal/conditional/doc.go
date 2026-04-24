// Package conditional provides a conditional processor for logpipe pipelines.
//
// A Processor wraps any transform function and executes it only when a
// field-level condition is satisfied. This allows selective enrichment,
// redaction, or mutation without writing separate routing rules.
//
// Supported operators:
//
//	"eq"     – field value equals the configured string
//	"neq"    – field value does not equal the configured string
//	"exists" – field is present in the JSON object (value is ignored)
//
// Example:
//
//	p, err := conditional.New("level", "eq", "error", redactor.Apply)
//	if err != nil { ... }
//	out, err := p.Apply(line)
package conditional
