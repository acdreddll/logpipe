// Package router implements log line fan-out routing for logpipe.
//
// A Router holds a set of named Routes. Each Route has an optional Match
// predicate and an io.Writer destination. When Dispatch is called with a
// raw JSON log line, the router evaluates each route's predicate and writes
// the line to every matching destination.
//
// Example usage:
//
//	r := router.New()
//	r.AddRoute(&router.Route{
//		Name: "stderr-errors",
//		Match: func(line []byte) bool {
//			return bytes.Contains(line, []byte(`"level":"error"`))
//		},
//		Writer: os.Stderr,
//	})
//	r.Dispatch(logLine)
package router
