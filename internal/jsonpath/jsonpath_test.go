package jsonpath_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/jsonpath"
)

func TestGet_TopLevel(t *testing.T) {
	line := []byte(`{"level":"info","msg":"hello"}`)
	v, err := jsonpath.Get(line, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "info" {
		t.Fatalf("expected info, got %v", v)
	}
}

func TestGet_Nested(t *testing.T) {
	line := []byte(`{"http":{"status":200}}`)
	v, err := jsonpath.Get(line, "http.status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.(float64) != 200 {
		t.Fatalf("expected 200, got %v", v)
	}
}

func TestGet_Missing(t *testing.T) {
	line := []byte(`{"level":"info"}`)
	_, err := jsonpath.Get(line, "missing.key")
	if err != jsonpath.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGet_InvalidJSON(t *testing.T) {
	_, err := jsonpath.Get([]byte(`not json`), "field")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSet_TopLevel(t *testing.T) {
	line := []byte(`{"level":"info"}`)
	out, err := jsonpath.Set(line, "level", "warn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]any
	json.Unmarshal(out, &m)
	if m["level"] != "warn" {
		t.Fatalf("expected warn, got %v", m["level"])
	}
}

func TestSet_CreatesNestedPath(t *testing.T) {
	line := []byte(`{"msg":"hi"}`)
	out, err := jsonpath.Set(line, "meta.env", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := jsonpath.Get(out, "meta.env")
	if err != nil || v != "prod" {
		t.Fatalf("expected prod, got %v (err=%v)", v, err)
	}
}

func TestDelete_TopLevel(t *testing.T) {
	line := []byte(`{"level":"info","msg":"hi"}`)
	out, err := jsonpath.Delete(line, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = jsonpath.Get(out, "level")
	if err != jsonpath.ErrNotFound {
		t.Fatal("expected field to be deleted")
	}
}

func TestDelete_MissingPath_NoError(t *testing.T) {
	line := []byte(`{"msg":"hi"}`)
	_, err := jsonpath.Delete(line, "does.not.exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
