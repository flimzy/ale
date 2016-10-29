package ale

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

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

type pageTemplate struct {
	filename string
	server   *Server
	lastRead time.Time
	template *template.Template
}

var pageTemplates = make(map[string]*pageTemplate)

const pageCacheTime = 15 * time.Second

func (s *Server) newPageTemplate(filename string) *pageTemplate {
	return &pageTemplate{
		filename: filename,
		server:   s,
	}
}

func (s *Server) getPageTemplate(name string) *pageTemplate {
	filename := s.TemplateDir + "/" + name
	p, ok := pageTemplates[filename]
	if !ok {
		p = s.newPageTemplate(filename)
		pageTemplates[filename] = p
	}
	return p
}

func (p *pageTemplate) read() error {
	if !p.shouldRead() {
		return nil
	}
	log.Printf("Reading template `%s`\n", p.filename)
	t := template.New("")
	if fm := p.server.View.FuncMap; fm != nil {
		t.Funcs(fm)
	}
	if _, err := t.ParseFiles(p.filename); err != nil {
		return errors.Wrapf(err, "Cannot read template '%s'", p.filename)
	}
	libPath := p.server.TemplateDir + "/lib"
	if _, err := os.Stat(libPath); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "Unable to read templates lib `%s`", libPath)
	} else if err == nil {
		_, err := t.ParseGlob(libPath + "/*")
		if err != nil {
			return errors.Wrapf(err, "Error parsing templates in '%s'", libPath)
		}
	}
	// Only set the template if there's no error. This allows us to keep using
	// an old template if an error is introduced.
	p.template = t
	p.lastRead = time.Now()
	return nil
}

func (p *pageTemplate) shouldRead() bool {
	if p.template == nil {
		return true
	}
	if time.Now().Sub(p.lastRead) < pageCacheTime {
		// Only check the mtime after dictionaryCacheTime
		return false
	}
	// I ignore errors here, as any failure will either re-surface
	// when we actually try to read the file, or it was transient
	info, _ := os.Stat(p.filename)
	if info.ModTime().After(p.lastRead) {
		return true
	}
	return false
}

func (s *Server) template(name string) (*template.Template, error) {
	if s.TemplateDir == "" {
		return nil, errors.Errorf("No TemplateDir specified")
	}
	t := s.getPageTemplate(name)
	if err := t.read(); err != nil {
		if t.template == nil {
			return nil, err
		}
		log.Printf("Error refreshing template '%s': %s\n", name, err.Error())
	}
	return t.template, nil
}

// Render renders the page
func (s *Server) Render(w ResponseWriter, r *http.Request) {
	if w.Written() {
		return
	}
	if err := s.renderTemplate(w, r); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
