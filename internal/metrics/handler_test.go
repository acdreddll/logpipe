package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	m := New()
	m.Received.Add(5)
	m.Emitted.Add(3)
	m.Dropped.Add(1)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	m.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("unexpected content-type: %s", ct)
	}

	var snap Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if snap.Received != 5 {
		t.Errorf("expected Received=5, got %d", snap.Received)
	}
	if snap.Emitted != 3 {
		t.Errorf("expected Emitted=3, got %d", snap.Emitted)
	}
	if snap.Dropped != 1 {
		t.Errorf("expected Dropped=1, got %d", snap.Dropped)
	}
}

func TestHandler_EmptyMetrics(t *testing.T) {
	m := New()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	m.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var snap Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if snap.Received != 0 || snap.Emitted != 0 || snap.Dropped != 0 {
		t.Errorf("expected zero snapshot, got %+v", snap)
	}
}
