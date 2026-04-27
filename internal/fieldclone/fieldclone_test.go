package fieldclone

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

func TestNew_EmptySource(t *testing.T) {
	_, err := New("", []string{"dst"})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNew_NoDestinations(t *testing.T) {
	_, err := New("src", nil)
	if err == nil {
		t.Fatal("expected error for empty destinations")
	}
}

func TestNew_EmptyDestinationName(t *testing.T) {
	_, err := New("src", []string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for empty destination name")
	}
}

func TestApply_ClonesFieldToSingleDest(t *testing.T) {
	c, err := New("level", []string{"severity"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := c.Apply([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Errorf("source field modified: got %v", m["level"])
	}
	if m["severity"] != "info" {
		t.Errorf("destination not set: got %v", m["severity"])
	}
}

func TestApply_ClonesFieldToMultipleDests(t *testing.T) {
	c, err := New("level", []string{"sev", "priority"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := c.Apply([]byte(`{"level":"warn"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["sev"] != "warn" {
		t.Errorf("sev not set: got %v", m["sev"])
	}
	if m["priority"] != "warn" {
		t.Errorf("priority not set: got %v", m["priority"])
	}
}

func TestApply_MissingSource_Unchanged(t *testing.T) {
	c, err := New("level", []string{"severity"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	input := []byte(`{"msg":"no level here"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, err := New("level", []string{"severity"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = c.Apply([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
