// Package ratelimit implements a token-bucket rate limiter for use in the
// logpipe pipeline. It allows callers to cap the number of log lines
// forwarded per second, protecting downstream outputs from bursts.
//
// Usage:
//
//	limiter, err := ratelimit.New(1000) // 1000 lines/sec
//	if err != nil { ... }
//	if limiter.Allow() {
//	    // forward the line
//	}
package ratelimit
