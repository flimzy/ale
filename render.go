package ale

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl = template.New("main")

func init() {
	_, err := tmpl.ParseGlob("html/*")
	if err != nil {
		log.Printf("Error parsing templates: %s\n", err)
		os.Exit(1)
	}
}

// Render renders the page
func (s *Server) Render(ctx context.Context, w ResponseWriter, r *http.Request) {
	if w.Written() {
		return
	}
	var template string
	if val := ctx.Value("template"); val != nil {
		template = val.(string)
	}
	//	fmt.Fprintf(w, "Found these templates: %s", tmpl.DefinedTemplates())
	if t := tmpl.Lookup(template); t != nil {
		t.Execute(w, ctx)
		return
	}
	fmt.Fprint(w, "Hello world!\n")
	fmt.Fprintf(w, "Requested asset = %s\n", r.URL.Path)
	fmt.Fprintf(w, "Requested template = %s\n", template)
}
