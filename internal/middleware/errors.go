package middleware

import (
	"garagesale/internal/platform/web"
	"log"
	"net/http"
)

func Errors(log *log.Logger) web.Middleware {
	// This is actual mw function to be executed
	f := func(before web.Handler) web.Handler {
		// This is main handler
		h := func(w http.ResponseWriter, h *http.Request) error {
			if err := before(w, h); err != nil {
				log.Printf("ERROR: %v", err)

				if err := web.RespondError(w, err); err != nil {
					return err
				}
			}

			return nil
		}

		return h
	}

	return f
}
