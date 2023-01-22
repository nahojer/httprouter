package httprouter

import "net/http"

// Middleware is a function designed to run some code before and/or after
// another [http.Handler].
type Middleware func(http.Handler) http.Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(mw []Middleware, h http.Handler) http.Handler {
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		wrap := mw[i]
		if wrap != nil {
			h = wrap(h)
		}
	}
	return h
}
