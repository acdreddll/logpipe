// Package truncate provides a Truncator that shortens string field values
// in JSON log lines exceeding a configured byte length.
//
// Example usage:
//
//	tr, err := truncate.New("message", 120, truncate.WithSuffix(" [truncated]"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := tr.Apply(line)
//
// Fields that are absent, non-string, or within the length limit are passed
// through unchanged. Invalid JSON returns an error.
package truncate
