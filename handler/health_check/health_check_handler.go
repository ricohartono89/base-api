package healthcheck

import (
	"net/http"

	"github.com/ricohartono89/base-api/db"
	"github.com/ricohartono89/base-api/errs"
	"github.com/ricohartono89/base-api/internal/handler"
	"github.com/ricohartono89/base-api/internal/json"
	"github.com/ricohartono89/base-api/services"
)

// Handler ...
type Handler struct {
	DB *db.DatabaseInterface
}

// NewHandler ...
func NewHandler(db *db.DatabaseInterface) Handler {
	return Handler{db}
}

// HealthCheckHandler ...
func (h Handler) HealthCheckHandler() handler.EndpointHandler {
	return func(w http.ResponseWriter, req *http.Request) (res json.HTTPResponse) {
		err := services.HealthCheck(*h.DB)
		if err != nil {
			return res.ImportJSONWebError(errs.BuildJSONWebError("1", "1", err, 500))
		}
		return res.SetOk(nil)
	}
}
