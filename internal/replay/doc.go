// Package replay provides controlled re-emission of stored log lines,
// allowing logpipe to replay historical log files through the pipeline at
// a configurable rate.
//
// A Replayer reads lines from any io.Reader and publishes them to a string
// channel. Callers can cap the number of lines emitted (WithMaxLines) or
// insert an inter-line delay to simulate live streaming (WithDelay).
//
// A Registry allows multiple named Replayer instances to be managed and
// retrieved by the rest of the application.
package replay
