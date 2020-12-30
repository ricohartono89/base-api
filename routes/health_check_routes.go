package routes

import (
	healthcheck "github.com/ricohartono89/base-api/handler/health_check"
)

func initHealthCheckEndpoints(p HandlerServiceProvider, r Context) {
	h := healthcheck.NewHandler(p.DB)

	r.RegisterEndpoint(getHealthCheckEndpoint(h))
}

func getHealthCheckEndpoint(h healthcheck.Handler) EndpointInfo {
	return EndpointInfo{
		HTTPMethod:    "GET",
		URLPattern:    "/health",
		Handler:       h.HealthCheckHandler(),
		Verifications: nil,
	}
}
