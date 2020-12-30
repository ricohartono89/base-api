package routes

import (
	"net/http"

	"github.com/ricohartono89/base-api/db"
	"github.com/ricohartono89/base-api/internal/handler"
	"github.com/ricohartono89/base-api/middleware"
)

// URLPrefix ...
const URLPrefix = "/api"

// HandlerServiceProvider ...
type HandlerServiceProvider struct {
	DB         *db.DatabaseInterface
	Redis      *db.RedisInterface
	HTTPClient *http.Client
}

// EndpointInfo ...
type EndpointInfo struct {
	HTTPMethod    string
	URLPattern    string
	Handler       handler.EndpointHandler
	Verifications []middleware.VerificationType
}

// InitRoutes ...
func InitRoutes(r Context, p HandlerServiceProvider) {
	initHealthCheckEndpoints(p, r)
}
