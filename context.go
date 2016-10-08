package ale

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/pkg/errors"
)

// ParamsContextKey is the key used to fetch URL paramaters from r.Context()
const ParamsContextKey = httptreemux.ParamsContextKey

// StashContextKey is the key used to fetch the stash from r.Context()
const StashContextKey = "🍺.stash"

// ViewContextKey is the key used to fetch the view from r.Context()
const ViewContextKey = "🍺.view"

// ClientIPContextKey is the key used to fetch the client IP from r.Context()
const ClientIPContextKey = "🍺.userip"

// ExtractClientIP parses the request and returns the clients IP
func ExtractClientIP(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, errors.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

// Context is a clone of the context.Context interface, for convenience
type Context interface {
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(interface{}) interface{}
}

// GetStash fetches the stash from an http.Request
func GetStash(r *http.Request) map[string]interface{} {
	stash, _ := r.Context().Value(StashContextKey).(map[string]interface{})
	return stash
}

// GetParams fetches the URL params from an http.Request
func GetParams(r *http.Request) map[string]string {
	params, _ := r.Context().Value(ParamsContextKey).(map[string]string)
	return params
}

// GetView fetches the view configuratoin from an http.Request
func GetView(r *http.Request) *View {
	view, _ := r.Context().Value(ViewContextKey).(*View)
	return view
}

// GetClientIP fetches the parsed client IP from an http.Request
func GetClientIP(r *http.Request) net.IP {
	ip, _ := r.Context().Value(ClientIPContextKey).(net.IP)
	return ip
}
