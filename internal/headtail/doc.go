// Package headtail implements a streaming head/tail processor for log lines.
//
// In Head mode the processor retains only the first N lines added since the
// last Flush, discarding any subsequent input until the buffer is drained.
//
// In Tail mode the processor maintains a sliding window of the most recent N
// lines, evicting the oldest entry whenever the limit is exceeded.
//
// Both modes are safe for concurrent use.
//
// Example (head):
//
//	p, _ := headtail.New(headtail.Head, 100)
//	for _, line := range lines {
//	    p.Add(line)
//	}
//	result := p.Flush()
package headtail
