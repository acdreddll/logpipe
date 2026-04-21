// Package flatten provides a Flattener that collapses nested JSON objects
// into a single-level map.
//
// Nested keys are joined with a configurable separator (default "."). Arrays
// are preserved as-is. An optional prefix can be prepended to every key in
// the output.
//
// Example
//
//	input:  {"http":{"method":"GET","status":200}}
//	output: {"http.method":"GET","http.status":200}
//
// A Registry is provided for managing multiple named Flattener instances
// within a single pipeline.
package flatten
