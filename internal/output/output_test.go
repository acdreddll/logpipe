package output_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/output"
)

func TestNew_Stdout(t *testing.T) {
	w, err := output.New("out", output.TypeStdout, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Type() != output.TypeStdout {
		t.Errorf("expected stdout type")
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := output.New("out", output.Type("kafka"), "")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNew_FileMissingPath(t *testing.T) {
	_, err := output.New("out", output.TypeFile, "")
	if err == nil {
		t.Fatal("expected error when file path is empty")
	}
}

func TestWrite_AppendsNewline(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "log.jsonl")
	w, err := output.New("file-out", output.TypeFile, tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = w.Write([]byte(`{"level":"info"}`))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !bytes.HasSuffix(data, []byte("\n")) {
		t.Errorf("expected newline at end, got: %q", data)
	}
}

func TestWrite_EmptyLine(t *testing.T) {
	w, err := output.New("err", output.TypeStderr, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, err := w.Write([]byte{})
	if err != nil || n != 0 {
		t.Errorf("expected no-op for empty line")
	}
}
