# `timeout` - HTTP Middleware

Provides timeout-related HTTP middleware for `go` HTTP servers.

The benefit of the `github.com/poly-gun/timeout` package includes built-in structured logging (`slog`),
and functional constructor patterns.

## Documentation

Official `godoc` documentation (with examples) can be found at the [Package Registry](https://pkg.go.dev/github.com/poly-gun/timeout).

## Usage

###### Add Package Dependency

```bash
go get -u github.com/poly-gun/timeout
```

###### Import & Implement

`main.go`

```go
package main

import (
    "math/rand"
    "net/http"
    "time"

    "github.com/poly-gun/timeout"
)

type Middlewares struct {
    middleware []func(http.Handler) http.Handler
}

func (m *Middlewares) Add(middlewares ...func(http.Handler) http.Handler) {
    if len(middlewares) == 0 {
        return
    }

    m.middleware = append(m.middleware, middlewares...)
}

func (m *Middlewares) Handler(parent http.Handler) (handler http.Handler) {
    var length = len(m.middleware)
    if length == 0 {
        return parent
    }

    // Wrap the end handler with the middleware chain
    handler = m.middleware[len(m.middleware)-1](parent)
    for i := len(m.middleware) - 2; i >= 0; i-- {
        handler = m.middleware[i](handler)
    }

    return
}

func Middleware() *Middlewares {
    return &Middlewares{
        middleware: make([]func(http.Handler) http.Handler, 0),
    }
}

func Example() {
    middlewares := Middleware()

    middlewares.Add(timeout.New().Options(func(o *timeout.Middleware) { o.Timeout = time.Second * 5 }).Handler)

    mux := http.NewServeMux()

    mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        process := time.Duration(rand.Intn(5)) * time.Second

        select {
        case <-ctx.Done():
            return

        case <-time.After(process):
            // The above channel simulates some hard work.
        }

        w.Write([]byte("done"))
    })

    http.ListenAndServe(":8080", middlewares.Handler(mux))
}
```

- Please refer to the [code examples](./example_test.go) for additional usage and implementation details.
- See https://pkg.go.dev/github.com/poly-gun/timeout for additional documentation.

## Contributions

See the [**Contributing Guide**](./CONTRIBUTING.md) for additional details on getting started.
