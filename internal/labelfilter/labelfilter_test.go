package labelfilter

import (
	"testing"
)

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", ModeAllow, []string{"info"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_UnknownMode(t *testing.T) {
	_, err := New("level", "bogus", []string{"info"})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNew_NoLabels(t *testing.T) {
	_, err := New("level", ModeAllow, nil)
	if err == nil {
		t.Fatal("expected error for empty labels")
	}
}

func TestKeep_Allow_MatchingLabel(t *testing.T) {
	f, _ := New("level", ModeAllow, []string{"info", "warn"})
	ok, err := f.Keep([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil || !ok {
		t.Fatalf("expected keep=true, got keep=%v err=%v", ok, err)
	}
}

func TestKeep_Allow_NonMatchingLabel(t *testing.T) {
	f, _ := New("level", ModeAllow, []string{"info"})
	ok, err := f.Keep([]byte(`{"level":"debug","msg":"verbose"}`))
	if err != nil || ok {
		t.Fatalf("expected keep=false, got keep=%v err=%v", ok, err)
	}
}

func TestKeep_Deny_MatchingLabel(t *testing.T) {
	f, _ := New("level", ModeDeny, []string{"debug"})
	ok, err := f.Keep([]byte(`{"level":"debug","msg":"verbose"}`))
	if err != nil || ok {
		t.Fatalf("expected keep=false, got keep=%v err=%v", ok, err)
	}
}

func TestKeep_Deny_NonMatchingLabel(t *testing.T) {
	f, _ := New("level", ModeDeny, []string{"debug"})
	ok, err := f.Keep([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil || !ok {
		t.Fatalf("expected keep=true, got keep=%v err=%v", ok, err)
	}
}

func TestKeep_FieldAbsent_Allow_Drops(t *testing.T) {
	f, _ := New("level", ModeAllow, []string{"info"})
	ok, err := f.Keep([]byte(`{"msg":"no level field"}`))
	if err != nil || ok {
		t.Fatalf("expected keep=false for absent field in allow mode, got keep=%v err=%v", ok, err)
	}
}

func TestKeep_FieldAbsent_Deny_Keeps(t *testing.T) {
	f, _ := New("level", ModeDeny, []string{"debug"})
	ok, err := f.Keep([]byte(`{"msg":"no level field"}`))
	if err != nil || !ok {
		t.Fatalf("expected keep=true for absent field in deny mode, got keep=%v err=%v", ok, err)
	}
}

func TestKeep_InvalidJSON_Dropped(t *testing.T) {
	f, _ := New("level", ModeAllow, []string{"info"})
	ok, err := f.Keep([]byte(`not json`))
	if err != nil || ok {
		t.Fatalf("expected keep=false for invalid JSON, got keep=%v err=%v", ok, err)
	}
}
