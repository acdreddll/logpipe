package fieldcount

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

func TestNew_DefaultField(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.field != defaultField {
		t.Errorf("expected default field %q, got %q", defaultField, c.field)
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New(WithField(""))
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestApply_InjectsCount(t *testing.T) {
	c, _ := New()
	out, err := c.Apply([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	v, ok := m[defaultField]
	if !ok {
		t.Fatalf("field %q not injected", defaultField)
	}
	// JSON numbers unmarshal as float64.
	if int(v.(float64)) != 3 {
		t.Errorf("expected count 3, got %v", v)
	}
}

func TestApply_CustomField(t *testing.T) {
	c, _ := New(WithField("num_fields"))
	out, err := c.Apply([]byte(`{"a":1,"b":2,"c":3,"d":4}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	v, ok := m["num_fields"]
	if !ok {
		t.Fatal("custom field not injected")
	}
	if int(v.(float64)) != 5 {
		t.Errorf("expected count 5, got %v", v)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := New()
	_, err := c.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_EmptyObject(t *testing.T) {
	c, _ := New()
	out, err := c.Apply([]byte(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	v := m[defaultField]
	// empty object: 0 original fields, then the count field itself is added (count was 0)
	if int(v.(float64)) != 0 {
		t.Errorf("expected count 0 for empty object, got %v", v)
	}
}
