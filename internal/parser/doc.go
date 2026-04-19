// Package parser implements log line parsing for logpipe.
//
// Supported formats:
//
//   - json:   standard JSON objects, one per line
//   - logfmt: key=value pairs as popularised by the logfmt library
//
// Usage:
//
//	p, err := parser.New("json")
//	if err != nil { ... }
//	m, err := p.Parse(line)
package parser
