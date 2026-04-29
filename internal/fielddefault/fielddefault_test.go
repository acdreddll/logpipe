package fielddefault

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", "val")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NilValue(t *testing.T) {
	_, err := New("env", nil)
	if err == nil {
		t.Fatal("expected error for nil value")
	}
}

func TestApply_FieldAbsent_InjectsDefault(t *testing.T) {
	d, err := New("env", "production")
	if err != nil {
		t.Fatal(err)
	}
	out, err := d.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["env"] != "production" {
		t.Fatalf("expected env=production, got %v", m["env"])
	}
}

func TestApply_FieldPresent_NotOverwritten(t *testing.T) {
	d, err := New("env", "production")
	if err != nil {
		t.Fatal(err)
	}
	out, err := d.Apply([]byte(`{"env":"staging"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["env"] != "staging" {
		t.Fatalf("expected env=staging, got %v", m["env"])
	}
}

func TestApply_WithOverwrite_ReplacesExisting(t *testing.T) {
	d, err := New("env", "production", WithOverwrite())
	if err != nil {
		t.Fatal(err)
	}
	out, err := d.Apply([]byte(`{"env":"staging"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["env"] != "production" {
		t.Fatalf("expected env=production, got %v", m["env"])
	}
}

func TestApply_EmptyStringField_InjectsDefault(t *testing.T) {
	d, err := New("env", "production")
	if err != nil {
		t.Fatal(err)
	}
	out, err := d.Apply([]byte(`{"env":""}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["env"] != "production" {
		t.Fatalf("expected env=production, got %v", m["env"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	d, err := New("env", "production")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
