package ale

import (
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/flimzy/log"
	"github.com/gorilla/handlers"
)

func (s *Server) logging(next http.Handler) http.Handler {
	if s.GetConf(ConfNoLog) != "" {
		return next
	}
	log.Debug("Enabling LoggingHandler\n")
	return handlers.LoggingHandler(os.Stderr, next)
}

func (s *Server) compress(next http.Handler) http.Handler {
	if s.GetConf(ConfNoCompress) != "" {
		return next
	}
	log.Debug("Enabling GzipHandler\n")
	return gziphandler.GzipHandler(next)
}
