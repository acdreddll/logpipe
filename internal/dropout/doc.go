// Package dropout provides a Dropper that suppresses structured log
// lines whose nominated field value matches a configured set of
// literal strings.
//
// Typical use-cases include:
//   - Silencing high-volume debug or trace lines in production.
//   - Discarding health-check or readiness-probe log noise before
//     routing to an expensive downstream sink.
//
// Usage:
//
//	d, err := dropout.New("level", []string{"debug", "trace"})
//	if err != nil { /* handle */ }
//
//	drop, err := d.ShouldDrop(line)
//	if err != nil { /* handle parse error */ }
//	if drop { continue }
//
// Multiple Dropper instances can be managed through a Registry.
package dropout
