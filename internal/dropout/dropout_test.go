package dropout

import (
	"testing"
)

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", []string{"debug"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NoValues(t *testing.T) {
	_, err := New("level", nil)
	if err == nil {
		t.Fatal("expected error for empty values")
	}
}

func TestShouldDrop_MatchingValue(t *testing.T) {
	d, err := New("level", []string{"debug", "trace"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := []byte(`{"level":"debug","msg":"verbose info"}`)
	drop, err := d.ShouldDrop(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drop {
		t.Fatal("expected line to be dropped")
	}
}

func TestShouldDrop_NonMatchingValue(t *testing.T) {
	d, _ := New("level", []string{"debug"})
	line := []byte(`{"level":"error","msg":"something broke"}`)
	drop, err := d.ShouldDrop(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drop {
		t.Fatal("expected line to be kept")
	}
}

func TestShouldDrop_MissingField(t *testing.T) {
	d, _ := New("level", []string{"debug"})
	line := []byte(`{"msg":"no level field"}`)
	drop, err := d.ShouldDrop(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drop {
		t.Fatal("expected line to be kept when field is absent")
	}
}

func TestShouldDrop_InvalidJSON(t *testing.T) {
	d, _ := New("level", []string{"debug"})
	_, err := d.ShouldDrop([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestShouldDrop_NonStringField(t *testing.T) {
	d, _ := New("code", []string{"200"})
	// field is a number, not a string — should not drop
	line := []byte(`{"code":200}`)
	drop, err := d.ShouldDrop(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drop {
		t.Fatal("expected non-string field to be kept")
	}
}
