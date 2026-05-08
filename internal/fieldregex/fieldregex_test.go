package fieldregex

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", `\d+`, "NUM")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("msg", "", "NUM")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("msg", `[invalid`, "X")
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestApply_ReplacesMatch(t *testing.T) {
	r, _ := New("msg", `\d+`, "NUM")
	out, err := r.Apply(`{"msg":"error 42 on line 7"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "error NUM on line NUM" {
		t.Fatalf("unexpected value: %v", m["msg"])
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	r, _ := New("msg", `\d+`, "NUM")
	input := `{"level":"info"}`
	out, err := r.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if out != input {
		t.Fatalf("expected unchanged, got %s", out)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	r, _ := New("count", `\d+`, "NUM")
	input := `{"count":99}`
	out, err := r.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["count"].(float64) != 99 {
		t.Fatalf("expected numeric field unchanged, got %v", m["count"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r, _ := New("msg", `\d+`, "NUM")
	_, err := r.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_EmptyReplacement(t *testing.T) {
	r, _ := New("msg", `\s+`, "")
	out, err := r.Apply(`{"msg":"hello world"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "helloworld" {
		t.Fatalf("unexpected value: %v", m["msg"])
	}
}
