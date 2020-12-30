package db

import "github.com/go-pg/pg/v9/orm"

// DatabaseInterface ...
type DatabaseInterface interface {
	Connect()
	Close()
	HealthCheck() error
	GetConnectionDB() orm.DB
	BeginTransaction() (DatabaseInterface, orm.DB, bool, error)
	Commit() error
	Rollback() error
	HealthCheckWithErrorReturnWhenDown() error
}
