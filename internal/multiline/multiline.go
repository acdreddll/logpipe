// Package multiline provides a log line assembler that merges continuation
// lines (e.g. Java stack traces) into a single JSON record.
package multiline

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// Assembler accumulates raw lines and flushes a merged JSON record once a
// "start" pattern is matched on the next incoming line.
type Assembler struct {
	startRe  *regexp.Regexp
	field    string
	pending  []string
}

// New creates an Assembler that treats any line matching startPattern as the
// beginning of a new log event. Continuation lines are appended to field
// (defaults to "message" when empty).
func New(startPattern, field string) (*Assembler, error) {
	if startPattern == "" {
		return nil, errors.New("multiline: startPattern must not be empty")
	}
	re, err := regexp.Compile(startPattern)
	if err != nil {
		return nil, err
	}
	if field == "" {
		field = "message"
	}
	return &Assembler{startRe: re, field: field}, nil
}

// Add feeds a raw line to the assembler. If the line starts a new event and
// there is a buffered event, the buffered event is returned as JSON together
// with a true flush flag. Otherwise (nil, false, nil) is returned.
func (a *Assembler) Add(line string) ([]byte, bool, error) {
	if a.startRe.MatchString(line) {
		if len(a.pending) == 0 {
			a.pending = append(a.pending, line)
			return nil, false, nil
		}
		out, err := a.flush()
		a.pending = []string{line}
		return out, true, err
	}
	a.pending = append(a.pending, line)
	return nil, false, nil
}

// Flush forces the current buffer to be emitted regardless of whether a new
// start line has been seen. Call this at end-of-stream.
func (a *Assembler) Flush() ([]byte, error) {
	if len(a.pending) == 0 {
		return nil, nil
	}
	return a.flush()
}

func (a *Assembler) flush() ([]byte, error) {
	merged := strings.Join(a.pending, "\n")
	a.pending = nil
	record := map[string]string{a.field: merged}
	return json.Marshal(record)
}
