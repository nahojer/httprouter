package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nahojer/httprouter"
)

var tests = []struct {
	RouteMethod  string
	RoutePattern string

	Method string
	Path   string
	Match  bool
	Params map[string]string
}{
	// prefix
	{
		"GET", "/not-prefix",
		"GET", "/not-prefix/anything/else", false, nil,
	},
	{
		"GET", "/prefixdots...",
		"GET", "/prefixdots/anything/else", true, nil,
	},
	{
		"GET", "/prefixdots/...",
		"GET", "/prefixdots", true, nil,
	},
	// path params
	{
		"GET", "/path-param/:id",
		"GET", "/path-param/123", true, map[string]string{"id": "123"},
	},
	{
		"GET", "/path-params/:era/:group/:member",
		"GET", "/path-params/60s/beatles/lennon", true, map[string]string{
			"era":    "60s",
			"group":  "beatles",
			"member": "lennon",
		},
	},
}

func TestRouter(t *testing.T) {

	// See example_test.go for middleware tests.

	for _, tt := range tests {
		r := httprouter.New()

		var (
			match    bool
			matchReq *http.Request
		)
		r.HandleFunc(tt.RouteMethod, tt.RoutePattern, func(w http.ResponseWriter, r *http.Request) {
			match = true
			matchReq = r
		})

		req := httptest.NewRequest(tt.Method, tt.Path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if match != tt.Match {
			t.Errorf("%q %q: got match %t, want %t", tt.Method, tt.Path, tt.Match, match)
		}

		for paramName, wantParamVal := range tt.Params {
			gotParamVal := httprouter.Param(matchReq, paramName)
			if gotParamVal != wantParamVal {
				t.Errorf("httprouter.Param(matchReq, %q) = %q, want %q", paramName, gotParamVal, wantParamVal)
			}
		}
	}
}
