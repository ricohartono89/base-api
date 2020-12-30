package handler

import (
	"github.com/gorilla/schema"
	"github.com/ricohartono89/base-api/db"
	"github.com/ricohartono89/base-api/domain"
)

// DefaultCurrentPage ...
const DefaultCurrentPage = 1

var decoder = schema.NewDecoder()
var defaultPageSize = 15

type Handler struct {
	DB    *db.DatabaseInterface
	Redis *db.RedisInterface
}

func (h *Handler) normPaging(i *domain.Paging) {
	if i.CurrentPage == 0 {
		i.CurrentPage = 1
	}

	if i.PageSize == 0 {
		i.PageSize = defaultPageSize
	}
}
