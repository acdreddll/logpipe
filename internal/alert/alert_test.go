package alert

import (
	"testing"
)

func noop(_, _ string) {}

func TestNew_Valid(t *testing.T) {
	_, err := New("test", Condition{Field: "level", Operator: "eq", Value: "error"}, noop)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidOperator(t *testing.T) {
	_, err := New("test", Condition{Field: "level", Operator: "gt"}, noop)
	if err == nil {
		t.Fatal("expected error for unknown operator")
	}
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("", Condition{Field: "level", Operator: "eq"}, noop)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestEvaluate_Eq_Fires(t *testing.T) {
	fired := false
	a, _ := New("err-alert", Condition{Field: "level", Operator: "eq", Value: "error"}, func(_, _ string) { fired = true })
	_ = a.Evaluate(`{"level":"error","msg":"oops"}`)
	if !fired {
		t.Fatal("expected alert to fire")
	}
}

func TestEvaluate_Eq_NoFire(t *testing.T) {
	fired := false
	a, _ := New("err-alert", Condition{Field: "level", Operator: "eq", Value: "error"}, func(_, _ string) { fired = true })
	_ = a.Evaluate(`{"level":"info","msg":"ok"}`)
	if fired {
		t.Fatal("expected alert not to fire")
	}
}

func TestEvaluate_Contains_Fires(t *testing.T) {
	fired := false
	a, _ := New("a", Condition{Field: "msg", Operator: "contains", Value: "fail"}, func(_, _ string) { fired = true })
	_ = a.Evaluate(`{"msg":"request failed"}`)
	if !fired {
		t.Fatal("expected alert to fire")
	}
}

func TestEvaluate_Exists_Fires(t *testing.T) {
	fired := false
	a, _ := New("a", Condition{Field: "error", Operator: "exists"}, func(_, _ string) { fired = true })
	_ = a.Evaluate(`{"error":"timeout"}`)
	if !fired {
		t.Fatal("expected alert to fire")
	}
}

func TestEvaluate_InvalidJSON(t *testing.T) {
	a, _ := New("a", Condition{Field: "level", Operator: "eq", Value: "error"}, noop)
	if err := a.Evaluate(`not-json`); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
