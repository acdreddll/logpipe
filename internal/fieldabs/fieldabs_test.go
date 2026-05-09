package fieldabs_test

import (
	"encoding/json"
	"testing"

	"logpipe/internal/fieldabs"
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
	_, err := fieldabs.New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_PositiveValue_Unchanged(t *testing.T) {
	p, _ := fieldabs.New("value")
	out, err := p.Apply(`{"value":42}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["value"].(float64) != 42 {
		t.Fatalf("expected 42, got %v", m["value"])
	}
}

func TestApply_NegativeValue_MadePositive(t *testing.T) {
	p, _ := fieldabs.New("value")
	out, err := p.Apply(`{"value":-7.5}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["value"].(float64) != 7.5 {
		t.Fatalf("expected 7.5, got %v", m["value"])
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	p, _ := fieldabs.New("missing")
	input := `{"other":1}`
	out, err := p.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if out != input {
		t.Fatalf("expected unchanged line, got %s", out)
	}
}

func TestApply_NonNumericField_Unchanged(t *testing.T) {
	p, _ := fieldabs.New("label")
	out, err := p.Apply(`{"label":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["label"].(string) != "hello" {
		t.Fatalf("expected hello, got %v", m["label"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := fieldabs.New("value")
	_, err := p.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_ZeroValue_Unchanged(t *testing.T) {
	p, _ := fieldabs.New("n")
	out, err := p.Apply(`{"n":0}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["n"].(float64) != 0 {
		t.Fatalf("expected 0, got %v", m["n"])
	}
}
