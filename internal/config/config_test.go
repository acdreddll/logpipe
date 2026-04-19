package config

import (
	"os"
	"testing"
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
	p := writeTemp(t, `
outputs:
  - name: stdout
    type: stdout
routes:
  - name: all
    output: stdout
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(cfg.Routes))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/no/such/file.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoRoutes(t *testing.T) {
	p := writeTemp(t, `outputs:\n  - name: x\n    type: stdout\n`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for no routes")
	}
}

func TestLoad_MissingOutputType(t *testing.T) {
	p := writeTemp(t, `
outputs:
  - name: bad
routes:
  - name: r
    output: bad
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing output type")
	}
}

func TestLoad_BufferConfig(t *testing.T) {
	p := writeTemp(t, `
outputs:
  - name: stdout
    type: stdout
routes:
  - name: buffered
    output: stdout
    buffer:
      size: 50
      interval: 1s
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Routes[0].Buffer == nil {
		t.Fatal("expected buffer config to be set")
	}
	if cfg.Routes[0].Buffer.Size != 50 {
		t.Fatalf("expected size 50, got %d", cfg.Routes[0].Buffer.Size)
	}
}
