// Package fieldtrim provides a Trimmer that removes leading and trailing
// whitespace (or a caller-supplied cutset) from a string field inside a
// JSON-encoded log event.
//
// # Usage
//
//	tr, err := fieldtrim.New("message")
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean, err := tr.Apply(line)
//
// A custom cutset can be supplied via WithCutset:
//
//	tr, _ := fieldtrim.New("tag", fieldtrim.WithCutset("-_"))
//
// Fields that are absent or not of type string are passed through without
// modification. Invalid JSON returns an error.
package fieldtrim
