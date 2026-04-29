// Package fielddefault provides a Defaulter that injects a configured value
// into a named JSON field when the field is absent or empty.
//
// # Usage
//
//	d, err := fielddefault.New("env", "production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := d.Apply(line)
//
// By default the field is only written when it is missing or an empty string.
// Use WithOverwrite to unconditionally set the field on every event.
package fielddefault
