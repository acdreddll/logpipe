// Package fieldsplit provides a processor that splits a delimited string field
// in a JSON log event into a JSON array.
//
// Example — split a comma-separated "tags" field:
//
//	s, err := fieldsplit.New("tags")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := s.Apply(`{"tags":"info,warning,debug"}`)
//	// out → {"tags":["info","warning","debug"]}
//
// Use WithSeparator to change the delimiter and WithDest to write the resulting
// array to a different field, leaving the original intact.
package fieldsplit
