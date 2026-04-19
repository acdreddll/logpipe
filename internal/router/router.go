package router

import (
	"io"
	"sync"
)

// Route defines a named output destination with an optional filter function.
type Route struct {
	Name   string
	Match  func(line []byte) bool
	Writer io.Writer
}

// Router fans out log lines to matching routes.
type Router struct {
	mu     sync.RWMutex
	routes []*Route
}

// New creates a new Router.
func New() *Router {
	return &Router{}
}

// AddRoute registers a new route.
func (r *Router) AddRoute(route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
}

// Dispatch sends the log line to all matching routes.
// Returns the number of routes the line was dispatched to.
func (r *Router) Dispatch(line []byte) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, route := range r.routes {
		if route.Match == nil || route.Match(line) {
			_, _ = route.Writer.Write(append(line, '\n'))
			count++
		}
	}
	return count
}

// Routes returns a copy of the current route names.
func (r *Router) Routes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, len(r.routes))
	for i, route := range r.routes {
		names[i] = route.Name
	}
	return names
}
