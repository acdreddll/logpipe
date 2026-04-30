// Package fieldcopy provides a log-line processor that copies the value of one
// JSON field into another field without removing the original.
//
// # Usage
//
//	c, err := fieldcopy.New("level", "severity")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := c.Apply(line)
//
// By default the destination field is NOT overwritten if it already exists.
// Use [WithOverwrite] to change this behaviour.
//
// # Registry
//
// Multiple named copiers can be managed through [Registry].
package fieldcopy
