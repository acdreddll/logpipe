package fieldexist

import (
	"testing"
)

func mustNew(t *testing.T, mode Mode, fields []string) *Filter {
	t.Helper()
	f, err := New(mode, fields)
	if err != nil {
		t.Fatalf("New: unexpected error: %v", err)
	}
	return f
}

func TestNew_NoFields(t *testing.T) {
	_, err := New(ModeRequire, nil)
	if err == nil {
		t.Fatal("expected error for empty fields, got nil")
	}
}

func TestNew_EmptyFieldName(t *testing.T) {
	_, err := New(ModeRequire, []string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for empty field name, got nil")
	}
}

func TestNew_UnknownMode(t *testing.T) {
	_, err := New(Mode(99), []string{"level"})
	if err == nil {
		t.Fatal("expected error for unknown mode, got nil")
	}
}

func TestKeep_Require_AllPresent(t *testing.T) {
	f := mustNew(t, ModeRequire, []string{"level", "msg"})
	ok, err := f.Keep([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Keep=true, got false")
	}
}

func TestKeep_Require_FieldMissing(t *testing.T) {
	f := mustNew(t, ModeRequire, []string{"level", "trace_id"})
	ok, err := f.Keep([]byte(`{"level":"info"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Keep=false, got true")
	}
}

func TestKeep_Exclude_NonePresent(t *testing.T) {
	f := mustNew(t, ModeExclude, []string{"debug", "trace_id"})
	ok, err := f.Keep([]byte(`{"level":"info","msg":"hi"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Keep=true, got false")
	}
}

func TestKeep_Exclude_FieldPresent(t *testing.T) {
	f := mustNew(t, ModeExclude, []string{"trace_id"})
	ok, err := f.Keep([]byte(`{"level":"info","trace_id":"abc"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Keep=false, got true")
	}
}

func TestKeep_InvalidJSON(t *testing.T) {
	f := mustNew(t, ModeRequire, []string{"level"})
	_, err := f.Keep([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
