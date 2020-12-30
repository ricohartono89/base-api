package main

import (
	"net/http"

	"github.com/ricohartono89/base-api/middleware"
	"github.com/ricohartono89/base-api/routes"
)

func (s *server) routes(client *http.Client, m middleware.Middleware) {
	router := routes.NewContextWithWay(s.router, s.client, routes.URLPrefix)

	provider := routes.HandlerServiceProvider{
		DB:         s.db,
		Redis:      s.redis,
		HTTPClient: s.client,
	}
	routes.InitRoutes(router, provider)
}
