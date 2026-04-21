package splitter

import (
	"testing"
)

func mustNew(t *testing.T, field string, routes map[string][]string, opts ...Option) *Splitter {
	t.Helper()
	s, err := New(field, routes, opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", map[string][]string{"a": {"b"}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyRoutes(t *testing.T) {
	_, err := New("level", map[string][]string{})
	if err == nil {
		t.Fatal("expected error for empty routes")
	}
}

func TestSplit_MatchedRoute(t *testing.T) {
	s := mustNew(t, "level", map[string][]string{
		"error": {"errors", "alerts"},
		"info":  {"general"},
	})

	buckets, err := s.Split([]byte(`{"level":"error","msg":"boom"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 2 || buckets[0] != "errors" || buckets[1] != "alerts" {
		t.Fatalf("unexpected buckets: %v", buckets)
	}
}

func TestSplit_NoMatchUsesDefault(t *testing.T) {
	s := mustNew(t, "level",
		map[string][]string{"error": {"errors"}},
		WithDefault("catchall"),
	)

	buckets, err := s.Split([]byte(`{"level":"debug","msg":"ok"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 1 || buckets[0] != "catchall" {
		t.Fatalf("expected catchall bucket, got %v", buckets)
	}
}

func TestSplit_MissingFieldUsesDefault(t *testing.T) {
	s := mustNew(t, "level",
		map[string][]string{"error": {"errors"}},
		WithDefault("catchall"),
	)

	buckets, err := s.Split([]byte(`{"msg":"no level here"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 1 || buckets[0] != "catchall" {
		t.Fatalf("expected catchall, got %v", buckets)
	}
}

func TestSplit_MissingFieldNoDefault_ReturnsNil(t *testing.T) {
	s := mustNew(t, "level", map[string][]string{"error": {"errors"}})

	buckets, err := s.Split([]byte(`{"msg":"no level"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 0 {
		t.Fatalf("expected nil/empty, got %v", buckets)
	}
}

func TestSplit_InvalidJSON(t *testing.T) {
	s := mustNew(t, "level", map[string][]string{"error": {"errors"}})
	_, err := s.Split([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestField(t *testing.T) {
	s := mustNew(t, "severity", map[string][]string{"high": {"pager"}})
	if s.Field() != "severity" {
		t.Fatalf("expected 'severity', got %q", s.Field())
	}
}
