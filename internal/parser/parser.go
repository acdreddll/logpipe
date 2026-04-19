// Package parser provides log line parsing utilities for logpipe.
package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format represents a supported log input format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatLogfmt Format = "logfmt"
)

// Parser parses raw log lines into JSON objects.
type Parser struct {
	format Format
}

// New returns a Parser for the given format string.
func New(format string) (*Parser, error) {
	f := Format(strings.ToLower(format))
	switch f {
	case FormatJSON, FormatLogfmt:
		return &Parser{format: f}, nil
	default:
		return nil, fmt.Errorf("parser: unsupported format %q", format)
	}
}

// Parse converts a raw log line into a map.
func (p *Parser) Parse(line string) (map[string]any, error) {
	switch p.format {
	case FormatJSON:
		return parseJSON(line)
	case FormatLogfmt:
		return parseLogfmt(line)
	default:
		return nil, fmt.Errorf("parser: unknown format %q", p.format)
	}
}

func parseJSON(line string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return nil, fmt.Errorf("parser: invalid JSON: %w", err)
	}
	return m, nil
}

func parseLogfmt(line string) (map[string]any, error) {
	m := make(map[string]any)
	for _, pair := range strings.Fields(line) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("parser: invalid logfmt pair %q", pair)
		}
		key := parts[0]
		val := strings.Trim(parts[1], `"`)
		m[key] = val
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("parser: empty logfmt line")
	}
	return m, nil
}
