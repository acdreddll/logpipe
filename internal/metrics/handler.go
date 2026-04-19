package metrics

import (
	"encoding/json"
	"net/http"
)

// Handler returns an HTTP handler that serves a JSON snapshot of m.
func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := m.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
		}
	}
}

// Serve starts a lightweight HTTP server on addr (e.g. ":9090") that exposes
// the /metrics endpoint. It blocks until the server returns an error.
func (m *Metrics) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", m.Handler())
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server.ListenAndServe()
}
