// Package enrichment provides field injection for structured log lines.
//
// An Enricher accepts a map of static key-value pairs and merges them into
// every JSON log line it processes. Existing keys in the log line are
// preserved and will not be overwritten by the enricher.
//
// Example usage:
//
//	enricher := enrichment.New(map[string]string{
//		"service": "my-service",
//		"env":     "production",
//	})
//	enriched, err := enricher.Enrich(logLine)
//
// Use NewRegistry to manage multiple named enrichers within a pipeline.
package enrichment
