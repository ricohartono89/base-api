package services

import (
	"github.com/ricohartono89/base-api/db"
)

// HealthCheck ...
func HealthCheck(db db.DatabaseInterface) error {
	return db.HealthCheck()
}
