// Package ratelimit implements a token-bucket rate limiter for use in the
// logpipe pipeline. It allows callers to cap the number of log lines
// forwarded per second, protecting downstream outputs from bursts.
//
// The rate limiter is safe for concurrent use by multiple goroutines.
//
// Usage:
//
//	limiter, err := ratelimit.New(1000) // 1000 lines/sec
//	if err != nil { ... }
//	if limiter.Allow() {
//	    // forward the line
//	}
//
// A rate of zero or negative is invalid and New will return an error.
// To disable rate limiting entirely, bypass the limiter rather than
// constructing one with a very large limit.
package ratelimit
