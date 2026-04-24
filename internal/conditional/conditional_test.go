package conditional

import (
	"encoding/json"
	"testing"
)

func identity(line []byte) ([]byte, error) { return line, nil }

func setFoo(line []byte) ([]byte, error) {
	var obj map[string]interface{}
	_ = json.Unmarshal(line, &obj)
	obj["transformed"] = true
	out, _ := json.Marshal(obj)
	return out, nil
}

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", "eq", "v", identity)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_UnknownOperator(t *testing.T) {
	_, err := New("level", "gte", "warn", identity)
	if err == nil {
		t.Fatal("expected error for unknown operator")
	}
}

func TestNew_NilApply(t *testing.T) {
	_, err := New("level", "eq", "error", nil)
	if err == nil {
		t.Fatal("expected error for nil apply")
	}
}

func TestApply_Eq_Matches(t *testing.T) {
	p, _ := New("level", "eq", "error", setFoo)
	out, err := p.Apply([]byte(`{"level":"error","msg":"boom"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["transformed"] != true {
		t.Error("expected transformed=true")
	}
}

func TestApply_Eq_NoMatch(t *testing.T) {
	p, _ := New("level", "eq", "error", setFoo)
	in := []byte(`{"level":"info","msg":"ok"}`)
	out, err := p.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(in) {
		t.Error("expected line unchanged")
	}
}

func TestApply_Exists_Matches(t *testing.T) {
	p, _ := New("trace_id", "exists", "", setFoo)
	out, _ := p.Apply([]byte(`{"trace_id":"abc","msg":"hi"}`))
	m := decode(t, out)
	if m["transformed"] != true {
		t.Error("expected transformed=true when field exists")
	}
}

func TestApply_Exists_NoMatch(t *testing.T) {
	p, _ := New("trace_id", "exists", "", setFoo)
	in := []byte(`{"msg":"hi"}`)
	out, _ := p.Apply(in)
	if string(out) != string(in) {
		t.Error("expected line unchanged when field absent")
	}
}

func TestApply_Neq_Matches(t *testing.T) {
	p, _ := New("level", "neq", "debug", setFoo)
	out, _ := p.Apply([]byte(`{"level":"warn"}`))
	m := decode(t, out)
	if m["transformed"] != true {
		t.Error("expected transformed=true for neq match")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := New("level", "eq", "error", identity)
	_, err := p.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
