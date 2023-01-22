package httprouter

import (
	"context"
	"net/http"

	"github.com/nahojer/sage"
)

// Router routes HTTP requests. Must be initialized by calling [New].
type Router struct {
	// NotFound is the [http.Handler] to call when no routes match. By default uses
	// the result of calling [http.NotFoundHandler].
	NotFound http.Handler
	// Middleware to run before and/or after all route handlers. These will be
	// executed in the order they are provided.
	Middleware []Middleware

	routes *sage.RoutesTrie[http.Handler]
}

// New returns a new Router initialized with default values.
func New() *Router {
	return &Router{
		NotFound: http.NotFoundHandler(),
		routes:   sage.NewRoutesTrie[http.Handler](),
	}
}

// Handle registers a [http.Handler] to run for a given HTTP method and path
// pair. Middleware wraps around the handler and are exectured in the order
// they are provided, and after any global middleware registered on the router.
//
// Path parameters are specified by prefixing path segments with a colon (":").
// To match any path that has a specific prefix, use the three dots ("...") prefix
// indicator. Examples:
//
//	// Call handleImages for any path prefixed with /images.
//	router.Handle("GET", "/images...", handleImages)
//
//	// Path parameter with name "id".
//	router.Handle("GET", "/users/:id", handleGetUser)
//
// Use [Param] to get the value of the path parameter from the request.
func (r *Router) Handle(method, pattern string, h http.Handler, mw ...Middleware) {
	h = wrapMiddleware(mw, h)
	h = wrapMiddleware(r.Middleware, h)
	r.routes.Add(method, pattern, h)
}

// HandleFunc is the [http.HandlerFunc] alternative to Handle.
func (r *Router) HandleFunc(method, pattern string, f http.HandlerFunc, mw ...Middleware) {
	r.Handle(method, pattern, f, mw...)
}

// ServeHTTP implements the [http.Handler] interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, params, found := r.routes.Lookup(req)
	if !found {
		r.NotFound.ServeHTTP(w, req)
		return
	}

	for k, v := range params {
		ctx := context.WithValue(req.Context(), ctxKey(k), v)
		req = req.WithContext(ctx)
	}

	h.ServeHTTP(w, req)
}
