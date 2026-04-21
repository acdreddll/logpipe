package schema

import (
	"testing"
)

func mustNew(t *testing.T, rules []FieldRule) *Validator {
	t.Helper()
	v, err := New(rules)
	if err != nil {
		t.Fatalf("New: unexpected error: %v", err)
	}
	return v
}

func TestNew_EmptyFieldName(t *testing.T) {
	_, err := New([]FieldRule{{Name: "", Required: true, Type: TypeString}})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New([]FieldRule{{Name: "level", Type: "object"}})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestValidate_RequiredFieldPresent(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "level", Required: true, Type: TypeString}})
	if err := v.Validate([]byte(`{"level":"info"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RequiredFieldMissing(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "level", Required: true, Type: TypeString}})
	if err := v.Validate([]byte(`{"msg":"hello"}`)); err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestValidate_OptionalFieldAbsent(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "trace_id", Required: false, Type: TypeString}})
	if err := v.Validate([]byte(`{"level":"info"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_WrongType(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "duration_ms", Required: true, Type: TypeNumber}})
	if err := v.Validate([]byte(`{"duration_ms":"fast"}`)); err == nil {
		t.Fatal("expected type mismatch error")
	}
}

func TestValidate_BooleanType(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "ok", Required: true, Type: TypeBoolean}})
	if err := v.Validate([]byte(`{"ok":true}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_AnyType_AcceptsAll(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "meta", Required: true, Type: TypeAny}})
	for _, line := range []string{`{"meta":42}`, `{"meta":"x"}`, `{"meta":true}`} {
		if err := v.Validate([]byte(line)); err != nil {
			t.Fatalf("unexpected error for %s: %v", line, err)
		}
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	v := mustNew(t, []FieldRule{{Name: "level", Required: true, Type: TypeString}})
	if err := v.Validate([]byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
