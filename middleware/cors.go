package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/ricohartono89/base-api/env"
	"github.com/thoas/go-funk"
)

func getAccessControlOriginVal(baseOrigin string) string {
	whitelistedOrigin := strings.Split(env.CorsOrigin(), ";")

	isWhitelistAll := funk.ContainsString(whitelistedOrigin, "*")
	if isWhitelistAll {
		return "*"
	}

	isWhitelisted := funk.ContainsString(whitelistedOrigin, baseOrigin)

	if isWhitelisted {
		return baseOrigin
	}

	return ""
}

// Cors is a middleware to handle Cross Origin Request
func (m Middleware) Cors(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseOrigin := r.Header.Get("Origin")

		if r.Method == "OPTIONS" {
			log.Print("preflight detected: ", r.Header)
			w.Header().Add("Connection", "keep-alive")
			w.Header().Add("Access-Control-Allow-Origin", getAccessControlOriginVal(baseOrigin))
			w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			w.Header().Add("Access-Control-Allow-Headers", "authorization, content-type, accept, accept-language")
			w.Header().Add("Access-Control-Max-Age", "86400")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", getAccessControlOriginVal(baseOrigin))
			h.ServeHTTP(w, r)
		}
	}
}
