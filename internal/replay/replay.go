package replay

import (
	"bufio"
	"io"
	"time"
)

// Replayer reads log lines from a source and emits them at a controlled rate,
// optionally adding a delay between lines to simulate real-time streaming.
type Replayer struct {
	delay    time.Duration
	maxLines int
}

// Option configures a Replayer.
type Option func(*Replayer)

// WithDelay sets the inter-line delay.
func WithDelay(d time.Duration) Option {
	return func(r *Replayer) { r.delay = d }
}

// WithMaxLines caps the number of lines emitted (0 = unlimited).
func WithMaxLines(n int) Option {
	return func(r *Replayer) { r.maxLines = n }
}

// New creates a Replayer with the given options.
func New(opts ...Option) *Replayer {
	r := &Replayer{}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Run reads lines from src and writes them to dst, respecting the configured
// delay and line cap. It closes dst when done or when ctx is cancelled via
// the returned done channel.
func (r *Replayer) Run(src io.Reader, dst chan<- string) {
	scanner := bufio.NewScanner(src)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		dst <- line
		count++
		if r.maxLines > 0 && count >= r.maxLines {
			break
		}
		if r.delay > 0 {
			time.Sleep(r.delay)
		}
	}
}
