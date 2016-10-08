package ale

import (
	"context"
	"net/http"
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

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	stash := make(map[string]interface{})
	stash["view"] = s.View.View
	stash["template"] = s.View.Template
	ctx := context.WithValue(s.Context, StashContextKey, stash)

	r := req.WithContext(ctx)
	w := &response{rw, false}
	s.router.ServeHTTP(w, r)
	s.Render(w, r)
}

// ServeFiles is a wrapper around httprouter.ServeFiles
func (s *Server) ServeFiles(path string, root http.FileSystem) {
	s.Router.GET(path, http.FileServer(root).ServeHTTP)
}
