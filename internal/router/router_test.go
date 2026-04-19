package router

import (
	"bytes"
	"testing"
)

func TestDispatch_AllRoutes(t *testing.T) {
	r := New()
	var buf1, buf2 bytes.Buffer

	r.AddRoute(&Route{Name: "all-1", Writer: &buf1})
	r.AddRoute(&Route{Name: "all-2", Writer: &buf2})

	line := []byte(`{"level":"info","msg":"hello"}`)
	count := r.Dispatch(line)

	if count != 2 {
		t.Fatalf("expected 2 dispatches, got %d", count)
	}
	if !bytes.Contains(buf1.Bytes(), line) {
		t.Error("buf1 missing line")
	}
	if !bytes.Contains(buf2.Bytes(), line) {
		t.Error("buf2 missing line")
	}
}

func TestDispatch_FilteredRoute(t *testing.T) {
	r := New()
	var buf bytes.Buffer

	r.AddRoute(&Route{
		Name: "errors-only",
		Match: func(line []byte) bool {
			return bytes.Contains(line, []byte(`"level":"error"`))
		},
		Writer: &buf,
	})

	r.Dispatch([]byte(`{"level":"info","msg":"ok"}`))
	if buf.Len() != 0 {
		t.Error("expected no output for info line")
	}

	errLine := []byte(`{"level":"error","msg":"fail"}`)
	r.Dispatch(errLine)
	if !bytes.Contains(buf.Bytes(), errLine) {
		t.Error("expected error line in output")
	}
}

func TestRoutes_Names(t *testing.T) {
	r := New()
	r.AddRoute(&Route{Name: "a", Writer: &bytes.Buffer{}})
	r.AddRoute(&Route{Name: "b", Writer: &bytes.Buffer{}})

	names := r.Routes()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected route names: %v", names)
	}
}
