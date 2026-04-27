package ceiling

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", 100)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_InfiniteMax(t *testing.T) {
	_, err := New("latency", 1.0/0.0)
	if err == nil {
		t.Fatal("expected error for infinite max")
	}
}

func TestApply_ClampsExceedingValue(t *testing.T) {
	c, _ := New("latency", 500.0)
	out, err := c.Apply([]byte(`{"latency":750.0}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["latency"].(float64) != 500.0 {
		t.Fatalf("expected 500, got %v", m["latency"])
	}
}

func TestApply_ValueBelowCeiling_Unchanged(t *testing.T) {
	c, _ := New("latency", 500.0)
	input := []byte(`{"latency":200.0}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["latency"].(float64) != 200.0 {
		t.Fatalf("expected 200, got %v", m["latency"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c, _ := New("latency", 500.0)
	input := []byte(`{"level":"info"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Fatalf("expected unchanged output, got %s", out)
	}
}

func TestApply_NonNumericField_Unchanged(t *testing.T) {
	c, _ := New("latency", 500.0)
	input := []byte(`{"latency":"fast"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["latency"].(string) != "fast" {
		t.Fatalf("expected 'fast', got %v", m["latency"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := New("latency", 500.0)
	_, err := c.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_ExactlyAtCeiling_Unchanged(t *testing.T) {
	c, _ := New("score", 100.0)
	input := []byte(`{"score":100.0}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["score"].(float64) != 100.0 {
		t.Fatalf("expected 100, got %v", m["score"])
	}
}
