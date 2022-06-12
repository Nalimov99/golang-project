package middleware

import (
	"expvar"
	"garagesale/internal/platform/web"
	"net/http"
	"runtime"
)

// m contains global program counters for the application
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("gorutines"),
	req: expvar.NewInt("reqests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counters for the application
func Metric() web.Middleware {
	// This is actual mw function to be executed
	f := func(before web.Handler) web.Handler {
		// This is main handler
		h := func(w http.ResponseWriter, r *http.Request) error {
			err := before(w, r)

			m.req.Add(1)

			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			if err != nil {
				m.err.Add(1)
			}

			// Return the error to be handled further up the chain
			return err
		}

		return h
	}

	return f
}
