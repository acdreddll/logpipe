package joinfield

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
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
	_, err := New("out", nil)
	if err == nil {
		t.Fatal("expected error for no sources")
	}
}

func TestNew_EmptySourceName(t *testing.T) {
	_, err := New("out", []string{"a", ""})
	if err == nil {
		t.Fatal("expected error for empty source name")
	}
}

func TestApply_JoinsFields(t *testing.T) {
	j, err := New("full", []string{"first", "last"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := j.Apply([]byte(`{"first":"Jane","last":"Doe"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["full"]; got != "Jane Doe" {
		t.Errorf("full = %q, want %q", got, "Jane Doe")
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	j, err := New("msg", []string{"a", "b", "c"}, WithSeparator("-"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := j.Apply([]byte(`{"a":"x","b":"y","c":"z"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["msg"]; got != "x-y-z" {
		t.Errorf("msg = %q, want %q", got, "x-y-z")
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	j, err := New("full", []string{"first", "middle", "last"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := j.Apply([]byte(`{"first":"Jane","last":"Doe"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["full"]; got != "Jane Doe" {
		t.Errorf("full = %q, want %q", got, "Jane Doe")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	j, err := New("out", []string{"a"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = j.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
