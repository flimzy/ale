package ale

import (
	"context"
	"html/template"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/flimzy/log"
	"github.com/tylerb/graceful"
)

// View provides view configuration options
type View struct {
	View     string
	Template string
	FuncMap  map[string]interface{}
}

// Timeout defines the default time to wait before killing active connections on shutdown or restart.
const Timeout = 10 * time.Second

// Server represents an Ale server instance.
type Server struct {
	// Timeout is the duration to wait before killing active requests when stopping the server
	Timeout time.Duration
	// Context is the master context for this server instance
	Context Context
	// TemplateDir is the name of the path which contains the HTML templates.
	// If this path contains a subdir 'lib/', any files contianed within lib
	// are also loaded into each template. This is where shared components
	// should generally go.
	TemplateDir string
	// View is the default View configuration
	View        *View
	templates   map[string]*template.Template
	router      *httptreemux.TreeMux
	Router      *httptreemux.ContextGroup
	httpServer  *graceful.Server
	httpsServer *graceful.Server
	envPrefix   string
	err         error
}

// New returns a new Ale server instance.
func New() *Server {
	router := httptreemux.New()
	s := &Server{
		Timeout: Timeout,
		Context: context.Background(),
		router:  router,
		Router:  router.UsingContext(),
	}
	return s
}

// Run initializes the web server instance
func (s *Server) Run() error {
	httpAddr := s.GetConf(ConfHTTPBind)
	httpsAddr := s.GetConf(ConfHTTPSBind)

	log.Debugf("Run(). httpAddr = %s, httpsAddr = %s", httpAddr, httpsAddr)

	if httpAddr != "" && httpsAddr != "" {
		return s.serveBoth(httpsAddr, httpAddr)
	}
	if httpAddr != "" {
		return s.serveHTTP(httpAddr)
	}
	return s.serveHTTPS(httpsAddr)
}
