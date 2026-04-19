// Package metrics provides lightweight, thread-safe atomic counters
// used throughout the logpipe pipeline to track throughput and error
// rates without external dependencies.
//
// Usage:
//
//	c := metrics.New()
//	c.LinesIn.Add(1)
//	snap := c.Snapshot()
//	fmt.Println(snap.LinesIn)
package metrics
