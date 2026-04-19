// Package output implements destination writers for logpipe.
//
// Supported output types:
//
//   - stdout: writes log lines to standard output
//   - stderr: writes log lines to standard error
//   - file:   appends log lines to a file at the specified path
//
// Usage:
//
//	w, err := output.New("my-output", output.TypeFile, "/var/log/app.jsonl")
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Write([]byte(`{"level":"error","msg":"oops"}`))
//
// Each Write call ensures the line ends with a newline character.
package output
