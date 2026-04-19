// Package pipeline wires together filter, router, and output into a
// single processing unit that reads log lines from an io.Reader.
package pipeline

import (
	"bufio"
	"io"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/router"
)

// Pipeline reads lines from a source, applies a filter, routes each
// matching line to the appropriate output via the router.
type Pipeline struct {
	filter *filter.Filter
	router *router.Router
	reg    *output.Registry
}

// New creates a Pipeline from the provided components.
func New(f *filter.Filter, r *router.Router, reg *output.Registry) *Pipeline {
	return &Pipeline{filter: f, router: r, reg: reg}
}

// Run reads lines from r until EOF or error, processing each one.
// It returns the first non-EOF read error encountered.
func (p *Pipeline) Run(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !p.filter.Match(line) {
			continue
		}
		if err := p.router.Dispatch(line, p.reg); err != nil {
			return err
		}
	}
	return scanner.Err()
}
