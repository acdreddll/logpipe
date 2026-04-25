// Package fingerprint provides deterministic SHA-256 fingerprinting for
// structured JSON log events.
//
// A Fingerprinter hashes the full event or a restricted set of fields,
// always sorting keys before hashing to ensure key-order independence.
//
// Typical use:
//
//	fp, err := fingerprint.New(fingerprint.WithFields("service", "level"))
//	if err != nil { ... }
//	hash, err := fp.Compute(line)
//
// The Registry allows multiple named Fingerprinter instances to be managed
// centrally and looked up by name at runtime.
package fingerprint
