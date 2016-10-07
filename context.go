package ale

import (
	"net/http"
	"time"
)

// ParamsContextKey is the key used to fetch URL paramaters from r.Context()
const ParamsContextKey = "🍺.params"

// StashContextKey is the key used to fetch the stash from r.Context()
const StashContextKey = "🍺.stash"

// Context is a clone of the context.Context interface, for convenience
type Context interface {
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(interface{}) interface{}
}

// Stash fetches the stash from an http.Request
func Stash(r *http.Request) map[string]interface{} {
	stash, _ := r.Context().Value(StashContextKey).(map[string]interface{})
	return stash
}

// Params fetches the URL params from an http.Request
func Params(r *http.Request) map[string]string {
	params, _ := r.Context().Value(ParamsContextKey).(map[string]string)
	return params
}
