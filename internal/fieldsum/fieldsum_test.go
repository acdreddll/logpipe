package fieldsum

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyDest(t *testing.T) {
	_, err := New("", []string{"a"})
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestNew_NoSources(t *testing.T) {
	_, err := New("total", nil)
	if err == nil {
		t.Fatal("expected error for no sources")
	}
}

func TestNew_EmptySourceName(t *testing.T) {
	_, err := New("total", []string{"a", ""})
	if err == nil {
		t.Fatal("expected error for empty source name")
	}
}

func TestApply_SumsFields(t *testing.T) {
	p, err := New("total", []string{"x", "y", "z"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := p.Apply([]byte(`{"x":1,"y":2,"z":3}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["total"].(float64) != 6 {
		t.Errorf("expected 6, got %v", m["total"])
	}
}

func TestApply_MissingFieldTreatedAsZero(t *testing.T) {
	p, _ := New("total", []string{"a", "b"})
	out, err := p.Apply([]byte(`{"a":5}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["total"].(float64) != 5 {
		t.Errorf("expected 5, got %v", m["total"])
	}
}

func TestApply_StringNumericField(t *testing.T) {
	p, _ := New("total", []string{"a", "b"})
	out, err := p.Apply([]byte(`{"a":"3.5","b":1.5}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["total"].(float64) != 5 {
		t.Errorf("expected 5, got %v", m["total"])
	}
}

func TestApply_NonNumericField_ReturnsError(t *testing.T) {
	p, _ := New("total", []string{"a"})
	_, err := p.Apply([]byte(`{"a":"not-a-number"}`))
	if err == nil {
		t.Fatal("expected error for non-numeric field")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := New("total", []string{"a"})
	_, err := p.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
