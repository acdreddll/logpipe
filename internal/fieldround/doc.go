// Package fieldround provides a processor that rounds a numeric JSON field
// to a configurable number of decimal places.
//
// # Usage
//
//	r, err := fieldround.New("latency_ms", 2)
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := r.Apply([]byte(`{"latency_ms":12.3456}`))
//	// out => {"latency_ms":12.35}
//
// Fields that are absent are passed through unchanged.
// Non-numeric fields return an error.
package fieldround
