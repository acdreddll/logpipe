package fieldboolean

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
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_AlreadyBool(t *testing.T) {
	c, _ := New("active")
	out, err := c.Apply([]byte(`{"active":true}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["active"] != true {
		t.Fatalf("expected true, got %v", decode(t, out)["active"])
	}
}

func TestApply_StringTrue(t *testing.T) {
	for _, s := range []string{"true", "1", "yes", "on", "TRUE", "Yes"} {
		c, _ := New("flag")
		out, err := c.Apply([]byte(`{"flag":"` + s + `"}`))
		if err != nil {
			t.Fatalf("string %q: unexpected error: %v", s, err)
		}
		if decode(t, out)["flag"] != true {
			t.Fatalf("string %q: expected true", s)
		}
	}
}

func TestApply_StringFalse(t *testing.T) {
	for _, s := range []string{"false", "0", "no", "off", "FALSE"} {
		c, _ := New("flag")
		out, err := c.Apply([]byte(`{"flag":"` + s + `"}`))
		if err != nil {
			t.Fatalf("string %q: unexpected error: %v", s, err)
		}
		if decode(t, out)["flag"] != false {
			t.Fatalf("string %q: expected false", s)
		}
	}
}

func TestApply_NumericNonZero(t *testing.T) {
	c, _ := New("n")
	out, err := c.Apply([]byte(`{"n":42}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["n"] != true {
		t.Fatal("expected true for non-zero number")
	}
}

func TestApply_NumericZero(t *testing.T) {
	c, _ := New("n")
	out, err := c.Apply([]byte(`{"n":0}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["n"] != false {
		t.Fatal("expected false for zero")
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	c, _ := New("missing")
	input := []byte(`{"other":"value"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Fatalf("expected unchanged output, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := New("flag")
	_, err := c.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_UnrecognisedString(t *testing.T) {
	c, _ := New("flag")
	_, err := c.Apply([]byte(`{"flag":"maybe"}`))
	if err == nil {
		t.Fatal("expected error for unrecognised string")
	}
}
