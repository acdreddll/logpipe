package truncate

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
	_, err := New("", 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroMaxLen(t *testing.T) {
	_, err := New("msg", 0)
	if err == nil {
		t.Fatal("expected error for zero maxLen")
	}
}

func TestApply_TruncatesLongValue(t *testing.T) {
	tr, err := New("msg", 5)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := tr.Apply([]byte(`{"msg":"hello world","level":"info"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "hello..." {
		t.Errorf("expected \"hello...\", got %q", m["msg"])
	}
}

func TestApply_ShortValueUnchanged(t *testing.T) {
	tr, _ := New("msg", 20)
	out, err := tr.Apply([]byte(`{"msg":"hi"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "hi" {
		t.Errorf("expected \"hi\", got %q", m["msg"])
	}
}

func TestApply_MissingFieldNoError(t *testing.T) {
	tr, _ := New("msg", 5)
	out, err := tr.Apply([]byte(`{"level":"info"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["msg"]; ok {
		t.Error("expected msg field to be absent")
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	tr, _ := New("count", 2)
	out, err := tr.Apply([]byte(`{"count":42}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["count"] != float64(42) {
		t.Errorf("expected 42, got %v", m["count"])
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	tr, _ := New("msg", 4, WithSuffix("~~"))
	out, err := tr.Apply([]byte(`{"msg":"abcdefgh"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "abcd~~" {
		t.Errorf("expected \"abcd~~\", got %q", m["msg"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr, _ := New("msg", 5)
	_, err := tr.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
