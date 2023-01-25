# httprouter

Idiomatic, fast and deliberately simple HTTP router for Go applications. Features:

* Path paramerization.
* Middleware, both router level and handler level.
* Prefix matching.
* Implements the `http.Handler` interface making it compliant with `http.ServeMux`.

Uses [sage](https://github.com/nahojer/sage) under the hood to match incoming HTTP requests to handlers.

All of the documentation can be found on the [go.dev](https://pkg.go.dev/github.com/nahojer/httprouter?tab=doc) website.

Is it Good? [Yes](https://news.ycombinator.com/item?id=3067434).
