// Package fieldregex provides a Replacer that applies a regular expression
// substitution to a single string field within a JSON log line.
//
// # Usage
//
//	r, err := fieldregex.New("message", `\d+`, "<NUM>")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := r.Apply(line)
//
// If the target field is absent or is not a string value the line is returned
// unchanged. All non-overlapping matches are replaced in a single pass using
// regexp.ReplaceAllString.
package fieldregex
