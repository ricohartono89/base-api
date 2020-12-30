package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ricohartono89/base-api/env"
)

type BearerToken = []byte

var jwtSigningKey = []byte(env.AuthJwtSigningKey())

var JwtTokenKey = env.AuthJwtTokenKey()

// VerificationType ...
type VerificationType int

// VerificationTypeConstants constant for verification type
var VerificationTypeConstants = struct {
	APIToken    VerificationType
	GoogleToken VerificationType
	JWTToken    VerificationType
	WorkerToken VerificationType
}{
	APIToken:    1,
	GoogleToken: 2,
	JWTToken:    3,
	WorkerToken: 4,
}

// Guest is a route middleware that will serve http handler
// if `Api Token` is valid.
// Use this middleware for routes that not require use to be logged
func (m Middleware) Guest(h http.Handler) http.HandlerFunc {
	return m.Group(h, false, m.ApiToken, m.JwtToken)
}

// Auth is a route middleware that will serve http handler
// if either `Api Token` or `Jwt Token` is valid
// Use this middleware for routes that require use to be logged.
func (m Middleware) Auth(h http.Handler) http.HandlerFunc {
	return m.Group(h, false, m.JwtToken)
}

// RetoolAuth is a route middleware that will handle Retool Authentication method
func (m Middleware) RetoolAuth(h http.Handler) http.HandlerFunc {
	return m.Group(h, false, m.ApiToken)
}

// ApiToken is a middleware to check Authorization Bearer Header
// is a valid `env` Api Token
func (m Middleware) ApiToken(w http.ResponseWriter, r *http.Request) (*http.Request, *Error) {
	token, err := GetBearerToken(r)

	if err != nil {
		return nil, &Error{err, http.StatusUnauthorized}
	}

	if env.APIToken() != string(token) {
		return nil, &Error{errors.New("invalid api token"), http.StatusUnauthorized}
	}

	return r, nil
}

// JwtToken is a middleware to check Authorization Bearer Header
// is a valid Jwt Token
func (m Middleware) JwtToken(w http.ResponseWriter, r *http.Request) (*http.Request, *Error) {
	token, err := GetBearerToken(r)

	if err != nil {
		return nil, &Error{err, http.StatusUnauthorized}
	}

	parsedToken, err := jwt.Parse(string(token), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSigningKey, nil
	})

	if err != nil {
		return nil, &Error{err, http.StatusUnauthorized}
	}

	requestWithContext := r.WithContext(context.WithValue(r.Context(), JwtTokenKey, parsedToken))

	return requestWithContext, nil
}

// GetBearerToken ...
func GetBearerToken(r *http.Request) (BearerToken, error) {
	authorizationHeader := r.Header.Get("Authorization")
	splitAuthorizationHeader := strings.Split(authorizationHeader, "Bearer")

	if len(splitAuthorizationHeader) != 2 {
		return nil, errors.New("invalid authorization bearer header")
	}

	token := strings.TrimSpace(splitAuthorizationHeader[1])

	return []byte(token), nil
}
