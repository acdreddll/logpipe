package fieldinvert

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_InvertsTrueToFalse(t *testing.T) {
	inv, err := New("active")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := inv.Apply(`{"active":true,"name":"svc"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["active"] != false {
		t.Errorf("expected active=false, got %v", m["active"])
	}
}

func TestApply_InvertsFalseToTrue(t *testing.T) {
	inv, _ := New("enabled")
	out, err := inv.Apply(`{"enabled":false}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", m["enabled"])
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	inv, _ := New("active")
	original := `{"name":"svc"}`
	out, err := inv.Apply(original)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out != original {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonBoolField_ReturnsError(t *testing.T) {
	inv, _ := New("active")
	_, err := inv.Apply(`{"active":"yes"}`)
	if err == nil {
		t.Fatal("expected error for non-boolean field")
	}
}

func TestApply_InvalidJSON_ReturnsError(t *testing.T) {
	inv, _ := New("active")
	_, err := inv.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
