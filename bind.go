package ale

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"sync"

	"github.com/flimzy/log"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/tylerb/graceful"
)

func (s *Server) mainHandler(addr string, shutdown func()) *graceful.Server {
	return &graceful.Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: s.logging(s.compress(s)),
		},
		Timeout:           s.Timeout,
		ShutdownInitiated: shutdown,
	}
}

// bindHTTPS binds to the HTTPS port or returns an error
func (s *Server) serveHTTPS(addr string) error {
	log.Printf("Binding HTTPS to %s", addr)
	s.httpsServer = s.mainHandler(addr, func() { log.Printf("Shutting down HTTPS service") })
	return s.httpsServer.ListenAndServeTLS(s.GetConf(ConfSSLCert), s.GetConf(ConfSSLKey))
}

// BindHTTP binds to the HTTP port or returns an error
func (s *Server) serveHTTP(addr string) error {
	return s.serveHTTPToHandler(addr, s)
}

func (s *Server) serveHTTPToHandler(addr string, h http.Handler) error {
	log.Printf("Binding HTTP to %s", addr)
	s.httpServer = s.mainHandler(addr, func() { log.Printf("Shutting down HTTP service") })
	return s.httpServer.ListenAndServe()
}

func (s *Server) serveFastCGI(addr string) error {
	log.Printf("Binding FastCGI to %s", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "Unable to find to FCGI bind point")
	}
	return fcgi.Serve(listener, s.logging(s.compress(s)))
	//	s.httpsServer = s.mainHandler(addr, func() { log.Printf("Shutting down FastCGI service") })
	//	return s.httpsServer.Serve(listener)
}

// serveBoth binds to the HTTP and HTTPS ports, redirecting all HTTP requests to HTTPS
func (s *Server) serveBoth(httpsAddr, httpAddr string) error {
	baseURI := s.GetConf("BASEURI")
	if baseURI == "" {
		return errors.Errorf("%s_BASEURI must be set to redirect from HTTP to HTTPS", s.envPrefix)
	}

	var httpsErr, httpErr error
	var wg *sync.WaitGroup
	wg.Add(1) // Yeah, we only wait for one to finish--then we kill the other
	go func() {
		defer wg.Done()
		httpsErr = s.serveHTTPS(httpsAddr)
	}()
	go func() {
		defer wg.Done()
		redirFunc := func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, baseURI, http.StatusFound)
		}
		httpErr = s.serveHTTPToHandler(httpAddr, http.HandlerFunc(redirFunc))
	}()
	wg.Wait()
	s.httpServer.Stop(s.Timeout)
	s.httpsServer.Stop(s.Timeout)
	<-s.httpServer.StopChan()
	<-s.httpsServer.StopChan()

	var err *multierror.Error
	if httpsErr != nil {
		multierror.Append(err, httpsErr)
	}
	if httpErr != nil {
		multierror.Append(err, httpErr)
	}
	return err
}
