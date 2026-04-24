// Package labelfilter implements a log-line gate that passes or discards
// structured JSON log entries based on whether a nominated field's value
// appears in a configured label set.
//
// Two modes are supported:
//
//	allow – only lines whose field value is in the label set are kept.
//	deny  – lines whose field value is in the label set are discarded.
//
// Example usage:
//
//	f, err := labelfilter.New("level", labelfilter.ModeAllow, []string{"info", "warn", "error"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	keep, _ := f.Keep(line)
//	if keep {
//		// forward line downstream
//	}
package labelfilter
