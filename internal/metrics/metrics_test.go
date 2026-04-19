package metrics_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.LinesIn != 0 || s.LinesFiltered != 0 || s.LinesRouted != 0 || s.LinesErrored != 0 {
		t.Fatalf("expected all zero counters, got %+v", s)
	}
}

func TestCounters_Increment(t *testing.T) {
	c := metrics.New()
	c.LinesIn.Add(5)
	c.LinesFiltered.Add(2)
	c.LinesRouted.Add(3)
	c.LinesErrored.Add(1)

	s := c.Snapshot()
	if s.LinesIn != 5 {
		t.Errorf("LinesIn: want 5, got %d", s.LinesIn)
	}
	if s.LinesFiltered != 2 {
		t.Errorf("LinesFiltered: want 2, got %d", s.LinesFiltered)
	}
	if s.LinesRouted != 3 {
		t.Errorf("LinesRouted: want 3, got %d", s.LinesRouted)
	}
	if s.LinesErrored != 1 {
		t.Errorf("LinesErrored: want 1, got %d", s.LinesErrored)
	}
}

func TestSnapshot_IsValueCopy(t *testing.T) {
	c := metrics.New()
	c.LinesIn.Add(10)
	s1 := c.Snapshot()
	c.LinesIn.Add(5)
	s2 := c.Snapshot()

	if s1.LinesIn != 10 {
		t.Errorf("s1.LinesIn should be 10, got %d", s1.LinesIn)
	}
	if s2.LinesIn != 15 {
		t.Errorf("s2.LinesIn should be 15, got %d", s2.LinesIn)
	}
}
