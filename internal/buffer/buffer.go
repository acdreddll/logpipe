// Package buffer provides a bounded, time-flushed ring buffer for batching
// log lines before they are dispatched downstream.
package buffer

import (
	"sync"
	"time"
)

// Buffer accumulates lines and flushes them either when the batch size is
// reached or when the flush interval elapses.
type Buffer struct {
	mu       sync.Mutex
	items    []string
	cap      int
	flushFn  func([]string)
	ticker   *time.Ticker
	stop     chan struct{}
}

// New creates a Buffer that calls flushFn with accumulated lines whenever
// size lines are buffered or interval elapses.
func New(size int, interval time.Duration, flushFn func([]string)) *Buffer {
	if size <= 0 {
		size = 100
	}
	b := &Buffer{
		items:   make([]string, 0, size),
		cap:     size,
		flushFn: flushFn,
		ticker:  time.NewTicker(interval),
		stop:    make(chan struct{}),
	}
	go b.run()
	return b
}

// Add enqueues a line. If the buffer is full it is flushed immediately.
func (b *Buffer) Add(line string) {
	b.mu.Lock()
	b.items = append(b.items, line)
	if len(b.items) >= b.cap {
		batch := b.drain()
		b.mu.Unlock()
		b.flushFn(batch)
		return
	}
	b.mu.Unlock()
}

// Stop flushes remaining items and stops the background ticker.
func (b *Buffer) Stop() {
	b.ticker.Stop()
	close(b.stop)
	b.mu.Lock()
	batch := b.drain()
	b.mu.Unlock()
	if len(batch) > 0 {
		b.flushFn(batch)
	}
}

func (b *Buffer) run() {
	for {
		select {
		case <-b.ticker.C:
			b.mu.Lock()
			batch := b.drain()
			b.mu.Unlock()
			if len(batch) > 0 {
				b.flushFn(batch)
			}
		case <-b.stop:
			return
		}
	}
}

func (b *Buffer) drain() []string {
	if len(b.items) == 0 {
		return nil
	}
	batch := b.items
	b.items = make([]string, 0, b.cap)
	return batch
}
