// Package metrics provides lightweight in-process counters for
// tracking log pipeline activity such as lines processed, filtered,
// routed, and errored.
package metrics

import "sync/atomic"

// Counters holds atomic counters for pipeline activity.
type Counters struct {
	LinesIn      atomic.Int64
	LinesFiltered atomic.Int64
	LinesRouted  atomic.Int64
	LinesErrored atomic.Int64
}

// New returns an initialised Counters instance.
func New() *Counters {
	return &Counters{}
}

// Snapshot returns a point-in-time copy of all counters as a plain struct.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		LinesIn:       c.LinesIn.Load(),
		LinesFiltered: c.LinesFiltered.Load(),
		LinesRouted:   c.LinesRouted.Load(),
		LinesErrored:  c.LinesErrored.Load(),
	}
}

// Snapshot is a value copy of Counters at a moment in time.
type Snapshot struct {
	LinesIn       int64
	LinesFiltered int64
	LinesRouted   int64
	LinesErrored  int64
}
