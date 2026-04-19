// Package output provides writers for routing log lines to various destinations.
package output

import (
	"fmt"
	"io"
	"os"
)

// Type represents the kind of output destination.
type Type string

const (
	TypeStdout Type = "stdout"
	TypeStderr Type = "stderr"
	TypeFile   Type = "file"
)

// Writer wraps an io.Writer with metadata.
type Writer struct {
	Name string
	t    Type
	w    io.Writer
}

// New creates a new Writer for the given type and optional path (used for file type).
func New(name string, t Type, path string) (*Writer, error) {
	switch t {
	case TypeStdout:
		return &Writer{Name: name, t: t, w: os.Stdout}, nil
	case TypeStderr:
		return &Writer{Name: name, t: t, w: os.Stderr}, nil
	case TypeFile:
		if path == "" {
			return nil, fmt.Errorf("output: file type requires a path")
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("output: failed to open file %q: %w", path, err)
		}
		return &Writer{Name: name, t: t, w: f}, nil
	default:
		return nil, fmt.Errorf("output: unknown type %q", t)
	}
}

// Write writes a log line (appends newline if missing).
func (w *Writer) Write(line []byte) (int, error) {
	if len(line) == 0 {
		return 0, nil
	}
	if line[len(line)-1] != '\n' {
		line = append(line, '\n')
	}
	return w.w.Write(line)
}

// Type returns the output type.
func (w *Writer) Type() Type {
	return w.t
}
