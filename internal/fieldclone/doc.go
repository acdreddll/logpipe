// Package fieldclone provides a Cloner that duplicates the value of a
// named JSON field into one or more destination fields without removing
// the original.
//
// Example usage:
//
//	c, err := fieldclone.New("level", []string{"severity", "priority"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := c.Apply([]byte(`{"level":"info","msg":"started"}`))
//	// out => {"level":"info","msg":"started","severity":"info","priority":"info"}
package fieldclone
