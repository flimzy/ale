package ale

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ResponseWriter is Ale's custom version of http.ResponseWriter
type ResponseWriter struct {
	r     http.ResponseWriter
	wrote bool
}

// Header returns the header map that will be (or was) sent by WriteHeader.
func (r *ResponseWriter) Header() http.Header {
	return r.r.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *ResponseWriter) Write(b []byte) (int, error) {
	r.wrote = true
	return r.r.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
func (r *ResponseWriter) WriteHeader(status int) {
	r.wrote = true
	r.r.WriteHeader(status)
}

// Written returns true if any data has already been sent to the client.
func (r *ResponseWriter) Written() bool {
	return r.wrote
}

// The Controller type is an adapter to allow the use of ordinary functions as HTTP handlers.
type Controller func(Context, http.ResponseWriter, *http.Request) Context

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Handle registers a new Controller with the given method and path.
func (s *Server) Handle(method, pattern string, c Controller) {
	s.router.Handle(method, pattern, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(s.Context, "params", p)
		s.Render(c(ctx, w, r), w, r)
	})
}

// Get is a shortcut for Handle("GET", pattern, c)
func (s *Server) Get(pattern string, c Controller) {
	s.Handle("GET", pattern, c)
}
