package routes

import (
	"net/http"
	"path"

	"github.com/ricohartono89/base-api/internal/router/way"
	"github.com/ricohartono89/base-api/middleware"
	"github.com/thoas/go-funk"
)

// Context ...
type Context struct {
	router     way.RouterInterface
	middleware *middleware.Middleware
	prefix     string
}

// NewContext ...
func NewContext(httpClient *http.Client, prefix string) Context {
	return Context{
		router:     way.NewRouter(),
		middleware: &middleware.Middleware{Client: httpClient},
		prefix:     prefix}
}

// NewContextWithWay ...
// Deprecated: this is a transition function to reuse 'way' created by server
func NewContextWithWay(router way.RouterInterface, httpClient *http.Client, prefix string) Context {
	return Context{
		router:     router,
		middleware: &middleware.Middleware{Client: httpClient},
		prefix:     prefix}
}

func (r *Context) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

// RegisterEndpoint ...
func (r *Context) RegisterEndpoint(info EndpointInfo) {
	r.RegisterEndpointWithPrefix(info, r.prefix)
}

// RegisterEndpointWithPrefix ...
func (r *Context) RegisterEndpointWithPrefix(info EndpointInfo, prefix string) {
	m := r.middleware
	urlPattern := getFullURLPattern(info, prefix)

	verificationFns := getVerificationMethod(m, info.Verifications)

	r.router.Handle(info.HTTPMethod, urlPattern, m.Cors(m.Verify(info.Handler, verificationFns...)))
}

func getVerificationMethod(m *middleware.Middleware, verifications []middleware.VerificationType) []middleware.MiddlewareFunc {
	return funk.Map(verifications, func(t middleware.VerificationType) middleware.MiddlewareFunc {
		switch t {
		case middleware.VerificationTypeConstants.JWTToken:
			return m.JwtToken
		default:
			return m.ApiToken
		}

	}).([]middleware.MiddlewareFunc)
}

func getFullURLPattern(info EndpointInfo, prefix string) string {
	if prefix == "" {
		return info.URLPattern
	}
	return path.Join(prefix, info.URLPattern)
}
