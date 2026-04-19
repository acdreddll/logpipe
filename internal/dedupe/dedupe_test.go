package dedupe

import (
	"testing"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := New(nil, 0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestIsDuplicate_FullLine(t *testing.T) {
	d, err := New(nil, 10)
	if err != nil {
		t.Fatal(err)
	}
	line := []byte(`{"level":"info","msg":"hello"}`)
	dup, err := d.IsDuplicate(line)
	if err != nil || dup {
		t.Fatalf("expected not duplicate, got dup=%v err=%v", dup, err)
	}
	dup, err = d.IsDuplicate(line)
	if err != nil || !dup {
		t.Fatalf("expected duplicate, got dup=%v err=%v", dup, err)
	}
}

func TestIsDuplicate_ByFields(t *testing.T) {
	d, err := New([]string{"msg"}, 10)
	if err != nil {
		t.Fatal(err)
	}
	a := []byte(`{"level":"info","msg":"hello"}`)
	b := []byte(`{"level":"warn","msg":"hello"}`)
	if dup, _ := d.IsDuplicate(a); dup {
		t.Fatal("first line should not be duplicate")
	}
	if dup, _ := d.IsDuplicate(b); !dup {
		t.Fatal("second line with same msg should be duplicate")
	}
}

func TestIsDuplicate_DifferentFields(t *testing.T) {
	d, err := New([]string{"msg"}, 10)
	if err != nil {
		t.Fatal(err)
	}
	a := []byte(`{"msg":"hello"}`)
	b := []byte(`{"msg":"world"}`)
	if dup, _ := d.IsDuplicate(a); dup {
		t.Fatal("first should not be dup")
	}
	if dup, _ := d.IsDuplicate(b); dup {
		t.Fatal("different msg should not be dup")
	}
}

func TestIsDuplicate_InvalidJSON(t *testing.T) {
	d, err := New([]string{"msg"}, 10)
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.IsDuplicate([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestIsDuplicate_CacheReset(t *testing.T) {
	d, err := New(nil, 2)
	if err != nil {
		t.Fatal(err)
	}
	lines := [][]byte{
		[]byte(`{"msg":"a"}`),
		[]byte(`{"msg":"b"}`),
		[]byte(`{"msg":"c"}`),
	}
	for _, l := range lines {
		d.IsDuplicate(l)
	}
	// after reset, first line should not be duplicate
	dup, err := d.IsDuplicate(lines[0])
	if err != nil {
		t.Fatal(err)
	}
	if dup {
		t.Fatal("expected cache reset; line should not be duplicate")
	}
}
