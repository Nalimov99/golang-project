package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ctxKey represents the type of value for the context key
type ctxKey int

// Key is used to store/retrieve Claims value fron context
const Key ctxKey = 1

// These are expected values for Claims.Roles
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// Claims represents the authorization claims transmited via a jwt
type Claims struct {
	Roles []string
	jwt.StandardClaims
}

// NewClaims construct a Claims value for the indetified user. The Claims
// expire within a specified duration of the provided time. Additional fields of the Claims
// can be set after calling NewClaims is desired
func NewClaims(subject string, roles []string, now time.Time, expires time.Duration) Claims {
	return Claims{
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expires).Unix(),
		},
	}
}
