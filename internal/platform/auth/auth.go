package auth

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// KeyLookupFunc is used to map a JWT key id to the corresponding public key.
// It is a requriement for creating an Authenticator
type KeyLookupFunc func(kid string) (*rsa.PublicKey, error)

// Authenticator is used to authenticate the clients. It can generate a token
// for a set of user claims and recreate the claims by passing the token
type Authenticator struct {
	privateKey           *rsa.PrivateKey
	activeKID            string
	algorithm            string
	publickKeyLookUpFunc KeyLookupFunc
	parser               *jwt.Parser
}

// NewSimpleKeyLookupFunc is a simple implementation of KeyFunc that only ever
// supports one key. This is easy for development but in production should be
// replaced with a caching layer that calls a JWKS endpoints
func NewSimpleKeyLookupFunc(activeKID string, publicKey *rsa.PublicKey) KeyLookupFunc {
	f := func(kid string) (*rsa.PublicKey, error) {
		if activeKID != kid {
			return nil, errors.Errorf("unrecognized key id %q", kid)
		}

		return publicKey, nil
	}

	return f
}

// NewAuthenticator creates an *Authenticator for use. It will error if:
//
// - the privateKey is nil
//
// - activeKID is blank
//
// - the algorithm is unsupported
//
// - publickLookUpFunc is not defined
func NewAuthenticator(
	privateKey *rsa.PrivateKey,
	acitiveKID, algorithm string,
	publickKeyLookUpFunc KeyLookupFunc,
) (*Authenticator, error) {

	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}

	if acitiveKID == "" {
		return nil, errors.New("acitiveKID cannot be blank")
	}

	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}

	if publickKeyLookUpFunc == nil {
		return nil, errors.New("publickKeyLookUpFunc key cannot be nil")
	}

	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := Authenticator{
		privateKey:           privateKey,
		activeKID:            acitiveKID,
		algorithm:            algorithm,
		publickKeyLookUpFunc: publickKeyLookUpFunc,
		parser:               &parser,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims
func (a *Authenticator) GenerateToken(claims Claims) (string, error) {
	method := jwt.GetSigningMethod(a.algorithm)
	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = a.activeKID

	str, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", errors.Wrapf(err, "signing token")
	}

	return str, nil
}
