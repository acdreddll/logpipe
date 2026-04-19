package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/router"
)

func newTestPipeline(t *testing.T, filterExpr string, routes []router.Route) (*pipeline.Pipeline, *output.Registry) {
	t.Helper()
	f, err := filter.New(filterExpr)
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}
	reg := output.NewRegistry()
	r := router.New(routes)
	return pipeline.New(f, r, reg), reg
}

func TestRun_FilterAndRoute(t *testing.T) {
	routes := []router.Route{
		{Name: "stdout", OutputType: "stdout"},
	}
	p, _ := newTestPipeline(t, `{"level":"error"}`, routes)

	input := strings.NewReader(
		`{"level":"error","msg":"boom"}` + "\n" +
			`{"level":"info","msg":"ok"}` + "\n" +
			`{"level":"error","msg":"oops"}` + "\n",
	)
	if err := p.Run(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	p, _ := newTestPipeline(t, `{}`, nil)
	if err := p.Run(strings.NewReader("")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_InvalidJSON_Skipped(t *testing.T) {
	p, _ := newTestPipeline(t, `{"level":"error"}`, nil)
	input := strings.NewReader("not-json\n{\"level\":\"error\"}\n")
	if err := p.Run(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
