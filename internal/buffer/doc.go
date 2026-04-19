// Package buffer provides a size- and time-bounded batch buffer for log lines.
//
// A Buffer accumulates incoming lines and forwards them as a slice to a
// caller-supplied flush function when either a configurable batch size is
// reached or a periodic interval elapses — whichever comes first.
//
// Typical usage:
//
//	b := buffer.New(200, 500*time.Millisecond, func(lines []string) {
//	    for _, l := range lines {
//	        router.Dispatch(l)
//	    }
//	})
//	defer b.Stop()
//
// Stop flushes any remaining lines and releases background resources.
package buffer
