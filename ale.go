package ale

import (
	"context"
	"html/template"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/flimzy/log"
	"github.com/pkg/errors"
	"github.com/tylerb/graceful"
)

// View provides view configuration options
type View struct {
	View     string
	Template string
	FuncMap  map[string]interface{}
}

// Copy makes a deep copy of a View
func (v *View) Copy() *View {
	nv := &View{
		View:     v.View,
		Template: v.Template,
		FuncMap:  make(map[string]interface{}),
	}
	for k, v := range v.FuncMap {
		nv.FuncMap[k] = v
	}
	return nv
}

var defaultFuncMap = map[string]interface{}{
	"makeMap": makeMap,
}

func makeMap(args ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for {
		if len(args) == 0 {
			return result
		}
		key, _ := args[0].(string)
		value := args[1]
		args = args[2:]
		result[key] = value
	}
}

// GetFuncMap returns a compiled FuncMap
func (v *View) GetFuncMap() template.FuncMap {
	fm := make(template.FuncMap)
	for k, v := range defaultFuncMap {
		fm[k] = v
	}
	for k, v := range v.FuncMap {
		fm[k] = v
	}
	return fm
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
	if s.GetConf(ConfFCGIBind) != "" {
		return s.FastCGI()
	}
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

// FastCGI binds to the FastCGI port
func (s *Server) FastCGI() error {
	fcgiAddr := s.GetConf(ConfFCGIBind)
	if fcgiAddr == "" {
		return errors.Errorf("%s_%s not set.", s.EnvPrefix(), ConfFCGIBind)
	}
	return s.serveFastCGI(fcgiAddr)
}
