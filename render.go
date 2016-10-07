package ale

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/flimzy/log"
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
	t, err := template.ParseFiles(tmplFile)
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
	stash := Stash(r)
	viewName, _ := stash["view"].(string)
	if viewName == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No view defined for %s.\n", r.URL.Path)
		return
	}
	tmplName, _ := stash["template"].(string)
	if tmplName == "" {
		tmplName = viewName
	}
	log.Debugf("viewName = %s, tmplName = %s\n", viewName, tmplName)
	if err := s.renderTemplate(w, viewName, tmplName, stash); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error executing template: %s\n", err)
	}
}

func (s *Server) renderTemplate(w ResponseWriter, view, tmpl string, stash map[string]interface{}) error {
	t, err := s.template(view)
	if err != nil {
		return err
	}
	buf := getBuf()
	defer putBuf(buf)
	log.Debugf("view = %s, tmpl = %s\n", view, tmpl)
	if err := t.ExecuteTemplate(buf, tmpl, stash); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
