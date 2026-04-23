package typecast

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

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", TypeInt)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New("field", Type("uuid"))
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestApply_StringToInt(t *testing.T) {
	c, _ := New("code", TypeInt)
	out, err := c.Apply([]byte(`{"code":"42","msg":"ok"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if v, ok := m["code"].(float64); !ok || v != 42 {
		t.Fatalf("expected code=42, got %v", m["code"])
	}
}

func TestApply_StringToFloat(t *testing.T) {
	c, _ := New("ratio", TypeFloat)
	out, err := c.Apply([]byte(`{"ratio":"3.14"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if v, ok := m["ratio"].(float64); !ok || v != 3.14 {
		t.Fatalf("expected ratio=3.14, got %v", m["ratio"])
	}
}

func TestApply_StringToBool(t *testing.T) {
	c, _ := New("active", TypeBool)
	out, err := c.Apply([]byte(`{"active":"true"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if v, ok := m["active"].(bool); !ok || !v {
		t.Fatalf("expected active=true, got %v", m["active"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c, _ := New("missing", TypeInt)
	input := []byte(`{"level":"info"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Fatalf("expected unchanged line, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := New("field", TypeInt)
	_, err := c.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_UncoercibleValue(t *testing.T) {
	c, _ := New("count", TypeInt)
	_, err := c.Apply([]byte(`{"count":"abc"}`))
	if err == nil {
		t.Fatal("expected error for unconvertible value")
	}
}
