package transform_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/transform"
)

func mustNew(t *testing.T, rules []transform.Rule) *transform.Transformer {
	t.Helper()
	tr, err := transform.New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return tr
}

func TestApply_Set(t *testing.T) {
	tr := mustNew(t, []transform.Rule{{Op: transform.OpSet, Field: "env", Value: "prod"}})
	out, err := tr.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(out), &m)
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", m["env"])
	}
}

func TestApply_Delete(t *testing.T) {
	tr := mustNew(t, []transform.Rule{{Op: transform.OpDelete, Field: "secret"}})
	out, err := tr.Apply(`{"level":"info","secret":"abc"}`)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(out), &m)
	if _, ok := m["secret"]; ok {
		t.Error("expected secret to be deleted")
	}
}

func TestApply_Rename(t *testing.T) {
	tr := mustNew(t, []transform.Rule{{Op: transform.OpRename, Field: "msg", To: "message"}})
	out, err := tr.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(out), &m)
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
	if _, ok := m["msg"]; ok {
		t.Error("old field should be removed")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr := mustNew(t, []transform.Rule{{Op: transform.OpSet, Field: "x", Value: "1"}})
	_, err := tr.Apply(`not-json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestNew_MissingField(t *testing.T) {
	_, err := transform.New([]transform.Rule{{Op: transform.OpSet}})
	if err == nil {
		t.Error("expected error when field is empty")
	}
}
