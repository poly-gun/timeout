package timeout

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

var key = "timeout"

type Middleware struct {
	// Timeout represents the duration to wait before considering an operation as timed out. If unspecified, or a negative value,
	// a default of 30 seconds is overwritten.
	Timeout time.Duration `json:"timeout"`
}

func (*Middleware) defaults() *Middleware {
	return &Middleware{
		Timeout: (time.Second * 30),
	}
}

func (m *Middleware) Options(options ...Variadic) *Middleware {
	if m == nil {
		*m = *(m.defaults())
	}

	for _, option := range options {
		option(m)
	}

	return m
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	if m.Timeout <= 0 {
		m.Timeout = (time.Second * 30)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		{
			value := m

			slog.Log(ctx, slog.LevelDebug, "Middleware", slog.Group("context", slog.String("key", string(key)), slog.Any("value", value)))

			ctx = context.WithValue(ctx, key, value)

			w.Header().Set("X-Timeout", m.Timeout.String())
		}

		ctx, cancel := context.WithTimeout(ctx, m.Timeout)
		defer func() {
			cancel()
			e := ctx.Err()
			if errors.Is(e, context.DeadlineExceeded) {
				http.Error(w, "gateway-timeout", http.StatusGatewayTimeout)
				return
			}
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Variadic represents a functional constructor for the [Middleware] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct a [Middleware] reference using packaged default(s).
type Variadic func(m *Middleware)

func New() *Middleware {
	return new(Middleware).defaults()
}
