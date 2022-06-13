package handlers

import (
	"context"
	"garagesale/internal/platform/auth"
	"garagesale/internal/platform/user"
	"garagesale/internal/platform/web"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// Users holds handlers for dealing with user
type Users struct {
	DB            *sqlx.DB
	Log           *log.Logger
	authenticator *auth.Authenticator
}

// Token generates an authentication token for a user. The client must include an email
// and password for the request using HTTP Basic Auth
func (u *Users) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.ContexValues)
	if !ok {
		return web.ErrContextValueMissing
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		return web.NewRequestError(errors.New("must provide email and password in Basic auth"), http.StatusUnauthorized)
	}

	claims, err := user.Authenticate(ctx, u.DB, v.Start, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "auth")
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}

	tkn.Token, err = u.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrapf(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
