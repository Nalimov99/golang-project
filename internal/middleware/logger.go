package middleware

import (
	"context"
	"errors"
	"garagesale/internal/platform/web"
	"log"
	"net/http"
	"time"
)

func Logger(log *log.Logger) web.Middleware {
	// This is actual mw function to be executed
	f := func(before web.Handler) web.Handler {
		// This is main handler
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := before(ctx, w, r)

			v, ok := ctx.Value(web.KeyValues).(*web.ContexValues)
			if !ok {
				return errors.New("web values missing from context")
			}

			log.Printf(
				"%d %s %s (%v)",
				v.StatusCode, r.Method, r.URL.Path, time.Since(v.Start),
			)

			// Return the error to be handled further up the chain
			return err
		}

		return h
	}

	return f
}
