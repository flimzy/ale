package ale

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ResponseWriter is Ale's custom version of http.ResponseWriter
type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(int)
	Written() bool
}

type response struct {
	rw    http.ResponseWriter
	wrote bool
}

// Header returns the header map that will be (or was) sent by WriteHeader.
func (r *response) Header() http.Header {
	return r.rw.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *response) Write(b []byte) (int, error) {
	r.wrote = true
	return r.rw.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
func (r *response) WriteHeader(status int) {
	r.wrote = true
	r.rw.WriteHeader(status)
}

// Written returns true if any data has already been sent to the client.
func (r *response) Written() bool {
	return r.wrote
}

// The Controller type is an adapter to allow the use of ordinary functions as HTTP handlers.
type Controller func(Context, http.ResponseWriter, *http.Request) Context

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Handle registers a new Controller with the given method and path.
func (s *Server) Handle(method, pattern string, c Controller, o ...RouteOptions) {
	s.router.Handle(method, pattern, func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(s.Context, "params", p)
		for _, opts := range o {
			if opts.Template != "" {
				ctx = context.WithValue(ctx, "template", opts.Template)
			}
		}
		w := &response{rw, false}
		s.Render(c(ctx, w, r), w, r)
	})
}

// Get is a shortcut for Handle("GET", pattern, c)
func (s *Server) Get(pattern string, c Controller, o ...RouteOptions) {
	s.Handle("GET", pattern, c, o...)
}

// ServeFiles is a wrapper around httprouter.ServeFiles
func (s *Server) ServeFiles(path string, root http.FileSystem) {
	s.router.ServeFiles(path, root)
}
