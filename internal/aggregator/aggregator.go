// Package aggregator provides field-based log aggregation with count, sum, and min/max support.
package aggregator

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Op defines the aggregation operation.
type Op string

const (
	OpCount Op = "count"
	OpSum   Op = "sum"
	OpMin   Op = "min"
	OpMax   Op = "max"
)

// Aggregator accumulates values for a given field using a specified operation.
type Aggregator struct {
	mu    sync.Mutex
	field string
	op    Op
	count int64
	sum   float64
	min   float64
	max   float64
	init  bool
}

// New creates a new Aggregator for the given field and operation.
func New(field string, op Op) (*Aggregator, error) {
	switch op {
	case OpCount, OpSum, OpMin, OpMax:
	default:
		return nil, fmt.Errorf("aggregator: unknown op %q", op)
	}
	if field == "" {
		return nil, fmt.Errorf("aggregator: field must not be empty")
	}
	return &Aggregator{field: field, op: op}, nil
}

// Add ingests a JSON log line and updates the aggregation state.
func (a *Aggregator) Add(line []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return fmt.Errorf("aggregator: invalid JSON: %w", err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.op == OpCount {
		a.count++
		return nil
	}
	raw, ok := m[a.field]
	if !ok {
		return nil
	}
	var v float64
	switch val := raw.(type) {
	case float64:
		v = val
	case int:
		v = float64(val)
	default:
		return fmt.Errorf("aggregator: field %q is not numeric", a.field)
	}
	a.count++
	a.sum += v
	if !a.init || v < a.min {
		a.min = v
	}
	if !a.init || v > a.max {
		a.max = v
	}
	a.init = true
	return nil
}

// Result returns the current aggregated value.
func (a *Aggregator) Result() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	switch a.op {
	case OpCount:
		return float64(a.count)
	case OpSum:
		return a.sum
	case OpMin:
		return a.min
	case OpMax:
		return a.max
	}
	return 0
}

// Reset clears the accumulated state.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.count = 0
	a.sum = 0
	a.min = 0
	a.max = 0
	a.init = false
}
