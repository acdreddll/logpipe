package fieldupper_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/fieldupper"
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
	_, err := fieldupper.New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_UppercasesField(t *testing.T) {
	tr, err := fieldupper.New("level")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := tr.Apply(`{"level":"warn","msg":"hello"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	m := decode(t, out)
	if got := m["level"]; got != "WARN" {
		t.Fatalf("expected WARN, got %v", got)
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	tr, _ := fieldupper.New("level")
	input := `{"msg":"hello"}`
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if decode(t, out)["msg"] != "hello" {
		t.Fatal("unexpected change")
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	tr, _ := fieldupper.New("code")
	input := `{"code":404}`
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if decode(t, out)["code"].(float64) != 404 {
		t.Fatal("numeric field should be unchanged")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr, _ := fieldupper.New("level")
	_, err := tr.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_AlreadyUppercase(t *testing.T) {
	tr, _ := fieldupper.New("level")
	out, err := tr.Apply(`{"level":"ERROR"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if decode(t, out)["level"] != "ERROR" {
		t.Fatal("already-uppercase value should stay ERROR")
	}
}
