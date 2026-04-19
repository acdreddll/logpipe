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
//
// Error handling:
//
// Run reads from the provided io.Reader line by line. Errors encountered
// while writing to an output are logged but do not halt processing; the
// pipeline continues with subsequent lines. A read error (other than
// io.EOF) causes Run to return that error immediately.
package pipeline
