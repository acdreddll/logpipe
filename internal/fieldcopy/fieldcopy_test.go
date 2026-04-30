package fieldcopy_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/fieldcopy"
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
	_, err := fieldcopy.New("", "dst")
	if err == nil {
		t.Fatal("expected error for empty src")
	}
}

func TestNew_EmptyDst(t *testing.T) {
	_, err := fieldcopy.New("src", "")
	if err == nil {
		t.Fatal("expected error for empty dst")
	}
}

func TestApply_CopiesField(t *testing.T) {
	c, _ := fieldcopy.New("level", "severity")
	out, err := c.Apply(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["severity"] != "info" {
		t.Errorf("expected severity=info, got %v", m["severity"])
	}
	if m["level"] != "info" {
		t.Errorf("source field should remain, got %v", m["level"])
	}
}

func TestApply_SourceAbsent_Unchanged(t *testing.T) {
	c, _ := fieldcopy.New("missing", "dst")
	original := `{"msg":"hello"}`
	out, err := c.Apply(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["dst"] != nil {
		t.Error("dst should not be set when source is absent")
	}
}

func TestApply_NoClobberByDefault(t *testing.T) {
	c, _ := fieldcopy.New("level", "severity")
	out, err := c.Apply(`{"level":"info","severity":"warn"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["severity"] != "warn" {
		t.Error("severity should not be overwritten without WithOverwrite")
	}
}

func TestApply_WithOverwrite(t *testing.T) {
	c, _ := fieldcopy.New("level", "severity", fieldcopy.WithOverwrite())
	out, err := c.Apply(`{"level":"info","severity":"warn"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["severity"] != "info" {
		t.Error("severity should be overwritten with WithOverwrite")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c, _ := fieldcopy.New("a", "b")
	out, err := c.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if out != `not-json` {
		t.Error("original line should be returned on error")
	}
}
