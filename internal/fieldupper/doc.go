// Package fieldupper provides a Transformer that converts the string value
// of a named JSON field to its uppercase equivalent.
//
// Fields that are absent, non-string, or already uppercase are passed through
// without modification. Non-JSON input is rejected with an error.
//
// Example
//
//	tr, err := fieldupper.New("level")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := tr.Apply(`{"level":"warn","msg":"disk full"}`)
//	// out => {"level":"WARN","msg":"disk full"}
package fieldupper
