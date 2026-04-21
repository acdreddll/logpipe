// Package schema validates structured JSON log lines against a declarative
// set of field rules.
//
// Each rule specifies a field name, whether the field is required, and the
// expected JSON value type (string, number, boolean, or any).
//
// Basic usage:
//
//	v, err := schema.New([]schema.FieldRule{
//		{Name: "level",  Required: true,  Type: schema.TypeString},
//		{Name: "ts",     Required: true,  Type: schema.TypeNumber},
//		{Name: "msg",    Required: false, Type: schema.TypeString},
//	})
//	if err != nil { /* handle */ }
//
//	if err := v.Validate(line); err != nil {
//		// line violates the schema
//	}
package schema
