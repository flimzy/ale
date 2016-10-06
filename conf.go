package ale

import (
	"github.com/flimzy/log"
	"github.com/spf13/viper"
)

const (
	// ConfHTTPBind is the config key for the HTTP bind address
	ConfHTTPBind = "HTTP_BIND"
	// ConfHTTPSBind is the config key for the HTTPS bind address
	ConfHTTPSBind = "HTTPS_BIND"
	// ConfSSLCert is the config key for the SSL Certificate location
	ConfSSLCert = "SSL_CERT"
	// ConfSSLKey is the config key for the SSL Key location
	ConfSSLKey = "SSL_KEY"
)

// Viper is an alias for viper.Viper
type Viper viper.Viper

func (s *Server) setConfDefaults() {
	log.Debug("Setting config defaults")
	s.viper.SetDefault("HTTP_BIND", ":8080")
	s.viper.BindEnv("HTTP_BIND")
}

// SetEnvPrefix sets the environment prefix for configuration
func (s *Server) SetEnvPrefix(prefix string) {
	log.Debugf("Setting ENV prefix to '%s'", prefix)
	s.envPrefix = prefix
	s.viper.SetEnvPrefix(prefix)
}

// SetConf sets the supplied configuration variable
func (s *Server) SetConf(key, value string) {
	s.viper.Set(key, value)
}

// GetConf retrieves the requested configuration variable
func (s *Server) GetConf(key string) string {
	value := s.viper.Get(key)
	log.Debugf("GetConf() %s = %s", key, value)
	if value == nil {
		return ""
	}
	return value.(string)
}
