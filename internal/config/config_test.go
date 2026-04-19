package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
routes:
  - name: stdout-all
    output:
      type: stdout
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(cfg.Routes))
	}
	if cfg.Routes[0].Name != "stdout-all" {
		t.Errorf("expected name stdout-all, got %q", cfg.Routes[0].Name)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoRoutes(t *testing.T) {
	path := writeTemp(t, `routes: []`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty routes")
	}
}

func TestLoad_MissingOutputType(t *testing.T) {
	path := writeTemp(t, `
routes:
  - name: bad-route
    output:
      path: /tmp/out.log
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing output type")
	}
}

func TestLoad_WithFilter(t *testing.T) {
	path := writeTemp(t, `
routes:
  - name: errors-only
    filter:
      level: error
    output:
      type: file
      path: /tmp/errors.log
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Routes[0].Filter["level"] != "error" {
		t.Errorf("expected filter level=error")
	}
}
