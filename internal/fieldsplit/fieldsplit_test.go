package fieldsplit_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/fieldsplit"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := fieldsplit.New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptySeparator(t *testing.T) {
	_, err := fieldsplit.New("tags", fieldsplit.WithSeparator(""))
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestNew_EmptyDest(t *testing.T) {
	_, err := fieldsplit.New("tags", fieldsplit.WithDest(""))
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestApply_SplitsCommaDelimited(t *testing.T) {
	s, _ := fieldsplit.New("tags")
	out, err := s.Apply(`{"tags":"a,b,c"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	arr, ok := m["tags"].([]any)
	if !ok || len(arr) != 3 {
		t.Fatalf("expected 3-element array, got %v", m["tags"])
	}
	if arr[1] != "b" {
		t.Errorf("expected arr[1]=b, got %v", arr[1])
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	s, _ := fieldsplit.New("path", fieldsplit.WithSeparator("/"))
	out, err := s.Apply(`{"path":"usr/local/bin"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	arr := m["path"].([]any)
	if len(arr) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(arr))
	}
}

func TestApply_WritesToDest(t *testing.T) {
	s, _ := fieldsplit.New("raw", fieldsplit.WithDest("parts"))
	out, err := s.Apply(`{"raw":"x,y"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["parts"]; !ok {
		t.Error("expected 'parts' key in output")
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	s, _ := fieldsplit.New("tags")
	const input = `{"level":"info"}`
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected unchanged line, got %s", out)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	s, _ := fieldsplit.New("count")
	const input = `{"count":42}`
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if v, _ := m["count"].(float64); v != 42 {
		t.Errorf("expected count=42, got %v", m["count"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s, _ := fieldsplit.New("tags")
	_, err := s.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
