// Package splitter provides field-based fan-out routing for structured log
// lines.
//
// A Splitter inspects a single JSON field and maps its value to one or more
// named output buckets. This allows a single log stream to be fanned out to
// multiple downstream consumers without duplicating filter logic.
//
// Example:
//
//	s, err := splitter.New("level", map[string][]string{
//		"error": {"errors", "pagerduty"},
//		"warn":  {"warnings"},
//	}, splitter.WithDefault("general"))
//
//	buckets, err := s.Split(line)
//
// Lines whose routing field is absent or whose value has no matching route
// are forwarded to the default bucket when one is configured, or silently
// dropped otherwise.
package splitter
