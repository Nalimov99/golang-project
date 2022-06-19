package middleware

import (
	"context"
	"errors"
	"garagesale/internal/platform/auth"
	"garagesale/internal/platform/web"
	"net/http"
	"strings"
)

// ErrForbidden is returned when an authenticated user
// does not have a sufficient role for an action
var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)

func Authenticate(authenticator *auth.Authenticator) web.Middleware {
	// This is actual mw function to be executed
	f := func(after web.Handler) web.Handler {
		// Wrap this handler around next provided
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("Expected Authorization header format: 'Bearer <token>'")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			claims, err := authenticator.ParseClaims(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = context.WithValue(ctx, auth.Key, claims)
			return after(ctx, w, r)
		}

		return h
	}

	return f
}

// HasRoles validates that an authenticated has at least one
// role from a specified list. This method constructs
// the actual function taht is used
func HasRoles(roles ...string) web.Middleware {
	// This is actual mw function to be executed
	f := func(after web.Handler) web.Handler {
		// Wrap this handler around next provided
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from context")
			}

			if !claims.HasRoles(roles...) {
				return ErrForbidden
			}

			return after(ctx, w, r)
		}

		return h
	}

	return f
}
