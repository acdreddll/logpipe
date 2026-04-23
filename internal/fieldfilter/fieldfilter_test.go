package fieldfilter

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

func TestNew_EmptyFields(t *testing.T) {
	_, err := New(ModeAllow, nil)
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_EmptyFieldName(t *testing.T) {
	_, err := New(ModeDeny, []string{""})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestApply_Allow_KeepsOnlyListed(t *testing.T) {
	ff, err := New(ModeAllow, []string{"level", "msg"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	input := []byte(`{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	out, err := ff.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to be present")
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected 'msg' to be present")
	}
	if _, ok := m["ts"]; ok {
		t.Error("expected 'ts' to be removed")
	}
}

func TestApply_Deny_RemovesListed(t *testing.T) {
	ff, err := New(ModeDeny, []string{"secret", "password"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	input := []byte(`{"level":"warn","secret":"abc","password":"xyz"}`)
	out, err := ff.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["secret"]; ok {
		t.Error("expected 'secret' to be removed")
	}
	if _, ok := m["password"]; ok {
		t.Error("expected 'password' to be removed")
	}
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to remain")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	ff, _ := New(ModeAllow, []string{"level"})
	_, err := ff.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_Allow_FieldAbsent_NoError(t *testing.T) {
	ff, _ := New(ModeAllow, []string{"missing"})
	out, err := ff.Apply([]byte(`{"level":"debug"}`))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if len(m) != 0 {
		t.Errorf("expected empty object, got %v", m)
	}
}
