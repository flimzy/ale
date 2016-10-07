package ale

import (
	"context"
	"fmt"
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Handle registers a new Controller with the given method and path.
func (s *Server) Handle(method, pattern string, h http.Handler, v ...View) {
	view := &View{}
	if s.View != nil {
		view.View = s.View.View
		view.Template = s.View.Template
	}
	for _, vv := range v {
		if vv.View != "" {
			view.View = vv.View
		}
		if vv.Template != "" {
			view.Template = vv.Template
		}
	}
	fmt.Printf("aggregate view = %v\n", view)
	s.router.Handle(method, pattern, func(rw http.ResponseWriter, req *http.Request, p map[string]string) {
		stash := map[string]interface{}{
			"view":     view.View,
			"template": view.Template,
		}
		ctx := context.WithValue(s.Context, StashContextKey, stash)
		ctx = context.WithValue(ctx, ParamsContextKey, p)
		w := &response{rw, false}
		r := req.WithContext(ctx)
		h.ServeHTTP(w, r)
		s.Render(w, r)
	})
}

// Get is a shortcut for Handle("GET", pattern, c)
func (s *Server) Get(pattern string, h http.Handler, v ...View) {
	s.Handle("GET", pattern, h, v...)
}

// ServeFiles is a wrapper around httprouter.ServeFiles
func (s *Server) ServeFiles(path string, root http.FileSystem) {
	s.Get(path, http.FileServer(root))
}
