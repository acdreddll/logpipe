package parser

import (
	"testing"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []string{"json", "logfmt", "JSON", "LOGFMT"} {
		_, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("csv")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestParse_JSON(t *testing.T) {
	p, _ := New("json")
	m, err := p.Parse(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	p, _ := New("json")
	_, err := p.Parse(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParse_Logfmt(t *testing.T) {
	p, _ := New("logfmt")
	m, err := p.Parse(`level=info msg="hello world"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
	if m["msg"] != "hello world" {
		t.Errorf("expected msg='hello world', got %v", m["msg"])
	}
}

func TestParse_Logfmt_InvalidPair(t *testing.T) {
	p, _ := New("logfmt")
	_, err := p.Parse(`levelinfo`)
	if err == nil {
		t.Fatal("expected error for invalid logfmt pair")
	}
}

func TestParse_Logfmt_EmptyLine(t *testing.T) {
	p, _ := New("logfmt")
	_, err := p.Parse("")
	if err == nil {
		t.Fatal("expected error for empty logfmt line")
	}
}
