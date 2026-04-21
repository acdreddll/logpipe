package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_MissingFile_OK(t *testing.T) {
	cp, err := New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s := cp.Get(); s.Offset != 0 || s.Source != "" {
		t.Fatalf("expected zero state, got %+v", s)
	}
}

func TestSave_And_Get(t *testing.T) {
	cp, _ := New(tmpPath(t))
	want := State{Offset: 42, Source: "app.log"}
	if err := cp.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got := cp.Get()
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestSave_Persists_Across_Reload(t *testing.T) {
	p := tmpPath(t)
	cp, _ := New(p)
	want := State{Offset: 99, Source: "svc.log"}
	cp.Save(want)

	cp2, err := New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := cp2.Get(); got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestReset_ClearsState(t *testing.T) {
	cp, _ := New(tmpPath(t))
	cp.Save(State{Offset: 7, Source: "x"})
	if err := cp.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if s := cp.Get(); s.Offset != 0 || s.Source != "" {
		t.Fatalf("state not cleared: %+v", s)
	}
}

func TestReset_RemovesFile(t *testing.T) {
	p := tmpPath(t)
	cp, _ := New(p)
	cp.Save(State{Offset: 5, Source: "y"})
	if err := cp.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("expected checkpoint file to be removed after Reset")
	}
}

func TestReset_NoFile_NoError(t *testing.T) {
	cp, _ := New(tmpPath(t))
	if err := cp.Reset(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	os.WriteFile(p, []byte("not json{"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for corrupt checkpoint file")
	}
}
