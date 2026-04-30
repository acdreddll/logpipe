package fieldmatch

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptySrc(t *testing.T) {
	_, err := New("", "matched", `\d+`)
	if err == nil {
		t.Fatal("expected error for empty src")
	}
}

func TestNew_EmptyDest(t *testing.T) {
	_, err := New("msg", "", `\d+`)
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("msg", "matched", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("msg", "matched", `[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestApply_MatchesPattern(t *testing.T) {
	m, err := New("msg", "matched", `error`)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := m.Apply(`{"msg":"an error occurred"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	obj := decode(t, out)
	if obj["matched"] != true {
		t.Errorf("expected matched=true, got %v", obj["matched"])
	}
}

func TestApply_NoMatch(t *testing.T) {
	m, _ := New("msg", "matched", `error`)
	out, err := m.Apply(`{"msg":"all good"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	obj := decode(t, out)
	if obj["matched"] != false {
		t.Errorf("expected matched=false, got %v", obj["matched"])
	}
}

func TestApply_MissingSrcField(t *testing.T) {
	m, _ := New("msg", "matched", `error`)
	out, err := m.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	obj := decode(t, out)
	if obj["matched"] != false {
		t.Errorf("expected matched=false for missing field, got %v", obj["matched"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m, _ := New("msg", "matched", `error`)
	_, err := m.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_NonStringField(t *testing.T) {
	m, _ := New("code", "matched", `42`)
	out, err := m.Apply(`{"code":42}`)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	// numeric field: no string match possible → false
	obj := decode(t, out)
	if obj["matched"] != false {
		t.Errorf("expected matched=false for non-string field, got %v", obj["matched"])
	}
}
