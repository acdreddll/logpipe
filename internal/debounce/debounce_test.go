package debounce

import (
	"testing"
	"time"
)

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", time.Second)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroWindow(t *testing.T) {
	_, err := New("msg", 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestAllow_FirstOccurrence(t *testing.T) {
	d, err := New("msg", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	line := []byte(`{"msg":"hello"}`)
	ok, err := d.Allow(line)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("first occurrence should be allowed")
	}
}

func TestAllow_SuppressedWithinWindow(t *testing.T) {
	d, err := New("msg", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	line := []byte(`{"msg":"hello"}`)
	d.Allow(line) //nolint
	ok, err := d.Allow(line)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("second occurrence within window should be suppressed")
	}
}

func TestAllow_AllowedAfterWindow(t *testing.T) {
	d, err := New("msg", 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	line := []byte(`{"msg":"hello"}`)
	d.Allow(line) //nolint
	time.Sleep(20 * time.Millisecond)
	ok, err := d.Allow(line)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("occurrence after window expiry should be allowed")
	}
}

func TestAllow_MissingField_AlwaysAllowed(t *testing.T) {
	d, err := New("msg", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	line := []byte(`{"level":"info"}`)
	for i := 0; i < 3; i++ {
		ok, err := d.Allow(line)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("line without key field should always be allowed (iteration %d)", i)
		}
	}
}

func TestAllow_InvalidJSON(t *testing.T) {
	d, err := New("msg", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.Allow([]byte(`not json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	d, err := New("msg", 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	d.Allow([]byte(`{"msg":"hello"}`)) //nolint
	time.Sleep(20 * time.Millisecond)
	d.Purge()
	d.mu.Lock()
	n := len(d.seen)
	d.mu.Unlock()
	if n != 0 {
		t.Errorf("expected 0 entries after purge, got %d", n)
	}
}
