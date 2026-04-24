// Package jsonpath provides dot-notation access to nested fields within
// JSON log lines.
//
// It exposes three operations:
//
//   - Get  – extract a value by path (e.g. "http.request.method")
//   - Set  – write a value at a path, creating intermediate maps as needed
//   - Delete – remove a field at a path
//
// All operations accept and return raw JSON byte slices so they compose
// naturally with the rest of the logpipe processing pipeline.
package jsonpath
