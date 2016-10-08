package ale

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/oxtoacart/bpool"
	"github.com/pkg/errors"
)

// BufPoolSize is the size of the BufferPool.
var BufPoolSize = 32

// BufPoolAlloc is the maximum size of each buffer.
var BufPoolAlloc = 10 * 1024

var bufpool *bpool.SizedBufferPool

func initBufPool() {
	bufpool = bpool.NewSizedBufferPool(BufPoolSize, BufPoolAlloc)
}

func getBuf() *bytes.Buffer {
	if bufpool == nil {
		initBufPool()
	}
	return bufpool.Get()
}

func putBuf(b *bytes.Buffer) {
	bufpool.Put(b)
}

func (s *Server) template(name string) (*template.Template, error) {
	if s.TemplateDir == "" {
		return nil, errors.Errorf("No TemplateDir specified")
	}
	tmplFile := s.TemplateDir + "/" + name
	if _, err := os.Stat(tmplFile); err != nil {
		return nil, errors.Wrapf(err, "Unable to read requested template `%s'", tmplFile)
	}
	t := template.New("")
	if s.View.FuncMap != nil {
		t = t.Funcs(s.View.FuncMap)
	}
	_, err := t.ParseFiles(tmplFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse template '%s'", tmplFile)
	}
	libPath := s.TemplateDir + "/lib"
	if _, err := os.Stat(libPath); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "Unable to read templates lib `%s`", libPath)
	} else if err == nil {
		_, err := t.ParseGlob(libPath + "/*")
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing templates in '%s'", libPath)
		}
	}
	return t, nil
}

// Render renders the page
func (s *Server) Render(w ResponseWriter, r *http.Request) {
	if w.Written() {
		return
	}
	if err := s.renderTemplate(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error executing template: %s\n", err)
	}
}

func (s *Server) renderTemplate(w ResponseWriter, r *http.Request) error {
	view := GetView(r)
	viewName := view.View
	if viewName == "" {
		return errors.Errorf("No view defined for %s", r.URL.Path)
	}
	tmplName := view.Template
	if tmplName == "" {
		tmplName = viewName
	}

	t, err := s.template(viewName)
	if err != nil {
		return err
	}
	if view.FuncMap != nil {
		t.Funcs(view.FuncMap)
	}
	buf := getBuf()
	defer putBuf(buf)
	if err := t.ExecuteTemplate(buf, tmplName, GetStash(r)); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
