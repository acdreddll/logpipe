// Package throttle provides cooldown-based suppression of repeated log lines.
//
// A Throttler tracks recently seen log lines and drops duplicates that arrive
// within a configurable cooldown window. Once the cooldown elapses the line
// is allowed through again.
//
// Lines are compared by their normalised JSON content so that key ordering
// differences in the raw bytes do not affect equality.
//
// Example usage:
//
//	th, err := throttle.New(5 * time.Second)
//	if err != nil { ... }
//	if ok, _ := th.Allow(line); ok {
//	    // forward line
//	}
package throttle
