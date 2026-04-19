// Package transform provides log record transformation rules for logpipe.
//
// A Transformer holds an ordered list of Rules. Each Rule specifies an
// operation (set, delete, or rename) and the target field. Rules are
// applied sequentially to every JSON log line that passes through the
// pipeline.
//
// Example usage:
//
//	tr, err := transform.New([]transform.Rule{
//		{Op: transform.OpSet,    Field: "env",     Value: "production"},
//		{Op: transform.OpDelete, Field: "password"},
//		{Op: transform.OpRename, Field: "msg",     To: "message"},
//	})
//	if err != nil { ... }
//	outLine, err := tr.Apply(inLine)
package transform
