// Package pipeline provides the top-level processing pipeline for logpipe.
//
// A Pipeline ties together three components:
//
//   - filter.Filter  – decides which log lines are relevant.
//   - router.Router  – maps a matching line to one or more named outputs.
//   - output.Registry – resolves output names to concrete io.Writers.
//
// Typical usage:
//
//	f, _ := filter.New(`{"level":"error"}`)
//	reg  := output.NewRegistry()
//	r    := router.New(routes)
//	p    := pipeline.New(f, r, reg)
//	p.Run(os.Stdin)
package pipeline
