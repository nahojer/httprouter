package httprouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/nahojer/httprouter"
)

func Example() {
	printer := func(text string) httprouter.Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(text)
				defer fmt.Println(text)
				next.ServeHTTP(w, r)
			})
		}
	}

	r := httprouter.New()
	r.Middleware = []httprouter.Middleware{printer("Global #1."), printer("Global #2.")}

	h := func(w http.ResponseWriter, r *http.Request) {
		band := httprouter.Param(r, "band")
		song := httprouter.Param(r, "song")
		fmt.Printf("Requested song %s by %s.\n", song, band)
	}
	r.HandleFunc("GET", "/music/:band/:song", h, printer("Local #1."), printer("Local #2."))

	req := httptest.NewRequest("GET", "http://localhost/music/U2/One", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Output:
	// Global #1.
	// Global #2.
	// Local #1.
	// Local #2.
	// Requested song One by U2.
	// Local #2.
	// Local #1.
	// Global #2.
	// Global #1.
}
