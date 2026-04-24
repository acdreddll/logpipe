// Package multiline implements a multi-line log assembler for logpipe.
//
// Many logging frameworks (Java, Python traceback, etc.) emit a single logical
// event across several physical lines. multiline.Assembler buffers raw lines
// and emits a single JSON record — with all continuation lines joined by \n
// inside a configurable field — whenever a new "start" pattern is detected or
// Flush is called at end-of-stream.
//
// Example usage:
//
//	asm, err := multiline.New(`^\d{4}-\d{2}-\d{2}`, "message")
//	for _, line := range lines {
//		if out, flushed, err := asm.Add(line); flushed {
//			process(out)
//		}
//	}
//	if out, _ := asm.Flush(); out != nil {
//		process(out)
//	}
package multiline
