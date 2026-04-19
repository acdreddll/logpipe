package filter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
)

func TestMatch_Eq(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "level", Operator: "eq", Value: "error"},
	})

	match, err := f.Match([]byte(`{"level":"error","msg":"oops"}`))
	if err != nil || !match {
		t.Fatalf("expected match, got match=%v err=%v", match, err)
	}

	match, err = f.Match([]byte(`{"level":"info","msg":"ok"}`))
	if err != nil || match {
		t.Fatalf("expected no match, got match=%v err=%v", match, err)
	}
}

func TestMatch_Contains(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "msg", Operator: "contains", Value: "timeout"},
	})

	match, _ := f.Match([]byte(`{"msg":"connection timeout reached"}`))
	if !match {
		t.Fatal("expected match for contains rule")
	}

	match, _ = f.Match([]byte(`{"msg":"all good"}`))
	if match {
		t.Fatal("expected no match")
	}
}

func TestMatch_Exists(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "trace_id", Operator: "exists"},
	})

	match, _ := f.Match([]byte(`{"trace_id":"abc123","msg":"traced"}`))
	if !match {
		t.Fatal("expected match when field exists")
	}

	match, _ = f.Match([]byte(`{"msg":"no trace"}`))
	if match {
		t.Fatal("expected no match when field absent")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f := filter.New([]filter.Rule{{Field: "level", Operator: "eq", Value: "error"}})
	_, err := f.Match([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
