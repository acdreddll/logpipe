// Package tagstrip provides a Stripper that removes named tag fields from
// structured JSON log lines.
//
// It is useful for scrubbing internal routing or diagnostic fields (e.g.
// "debug", "trace", "_internal") before forwarding log events downstream.
//
// Usage:
//
//	s, err := tagstrip.New([]string{"debug", "_routing"})
//	if err != nil { ... }
//	clean, err := s.Apply(rawLine)
//
// A Registry is provided for managing multiple named Stripper instances,
// which can be looked up by name during pipeline configuration.
package tagstrip
