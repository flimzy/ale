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
	ctx := s.Context
	ip, err := ExtractClientIP(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctx = context.WithValue(ctx, ClientIPContextKey, ip)

	stash := make(map[string]interface{})
	stash["view"] = s.View.View
	stash["template"] = s.View.Template
	stash["req"] = req
	ctx = context.WithValue(ctx, StashContextKey, stash)
	ctx = context.WithValue(ctx, ViewContextKey, s.View.Copy())

	r := req.WithContext(ctx)
	w := &response{rw, false}
	s.router.ServeHTTP(w, r)
	s.Render(w, r)
}

// ServeFiles is a wrapper around httprouter.ServeFiles
func (s *Server) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	s.Router.GET(path, func(w http.ResponseWriter, req *http.Request) {
		params := GetParams(req)
		req.URL.Path = params["filepath"]
		fileServer.ServeHTTP(w, req)
	})
}
