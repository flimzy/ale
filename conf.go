package ale

import (
	"os"

	"github.com/flimzy/log"
)

var defaults = map[string]string{
	"HTTP_BIND": ":8080",
}

const (
	// ConfHTTPBind is the config key for the HTTP bind address
	ConfHTTPBind = "HTTP_BIND"
	// ConfHTTPSBind is the config key for the HTTPS bind address
	ConfHTTPSBind = "HTTPS_BIND"
	// ConfSSLCert is the config key for the SSL Certificate location
	ConfSSLCert = "SSL_CERT"
	// ConfSSLKey is the config key for the SSL Key location
	ConfSSLKey = "SSL_KEY"
	// ConfFCGIBind will enable FastCGI mode if set to a bind address
	ConfFCGIBind = "FASTCGI_BIND"
)

// SetEnvPrefix sets the environment prefix for configuration
func (s *Server) SetEnvPrefix(prefix string) {
	log.Debugf("Setting ENV prefix to '%s'", prefix)
	s.envPrefix = prefix
}

// EnvPrefix returns the currently configured ENV prefix
func (s *Server) EnvPrefix() string {
	return s.envPrefix
}

// GetConf retrieves the requested configuration variable
func (s *Server) GetConf(key string) (val string) {
	val, _ = defaults[key] // Set default, if there is one
	if s.envPrefix != "" {
		return os.Getenv(s.envPrefix + "_" + key)
	}
	return os.Getenv(key)
}
