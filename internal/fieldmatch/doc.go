// Package fieldmatch evaluates a regular expression against a named field in
// a structured JSON log line and writes the boolean result (true / false) into
// a configurable destination field.
//
// Example usage:
//
//	// Mark lines whose "msg" field contains the word "error".
//	 m, err := fieldmatch.New("msg", "is_error", `(?i)error`)
//	 if err != nil {
//	     log.Fatal(err)
//	 }
//	 out, err := m.Apply(line)
//
// The destination field is always written, even when the source field is
// absent or is not a string — in those cases the result is false.
package fieldmatch
