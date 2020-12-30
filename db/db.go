package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/ricohartono89/base-api/env"
)

type debugHook struct{}

var _ pg.QueryHook = (*debugHook)(nil)

// BeforeQuery ...
func (debugHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		log.Printf("Error %s executing query:\n%s\n", evt.Err, q)
	} else {
		log.Printf("%s", q)
	}

	return ctx, nil
}

// AfterQuery
func (debugHook) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}

// Database ...
type Database struct {
	innerDB       *pg.DB
	transactionDB *pg.Tx
}

// HealthCheck ...
func (db *Database) HealthCheck() error {
	_, err := db.innerDB.Exec("SELECT 1")
	if err != nil {
		fmt.Println("PostgreSQL is down")
	}

	return nil
}

// Connect ...
func (db *Database) Connect() {
	var tlsConfig *tls.Config

	if env.Env() == "production" || env.Env() == "staging" || env.Env() == "development" {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	db.innerDB = pg.Connect(&pg.Options{
		ApplicationName: env.AppName(),
		User:            env.DBUser(),
		Password:        env.DBPassword(),
		Addr:            env.DBAddress(),
		Database:        env.DBDatabase(),
		TLSConfig:       tlsConfig,
	})

	if env.Env() == "local" {
		db.innerDB.AddQueryHook(debugHook{})
	}
}

// Close ...
func (db *Database) Close() {
	db.innerDB.Close()
}

// GetConnectionDB ...
func (db *Database) GetConnectionDB() orm.DB {
	if db.transactionDB != nil {
		return db.transactionDB
	}

	return db.innerDB
}

// BeginTransaction ...
func (db *Database) BeginTransaction() (DatabaseInterface, orm.DB, bool, error) {
	if db.transactionDB != nil {
		return db, db.transactionDB, true, nil
	}

	newDB := Database{innerDB: db.innerDB}
	tx, err := newDB.innerDB.Begin()
	if err != nil {
		return nil, nil, false, err
	}

	newDB.transactionDB = tx
	return &newDB, newDB.transactionDB, false, nil
}

// Commit ...
func (db *Database) Commit() error {
	if db.transactionDB == nil {
		return nil
	}

	err := db.transactionDB.Commit()
	if err != nil {
		return db.transactionDB.Close()
	}

	return db.closeTransactionDB()
}

func (db *Database) closeTransactionDB() error {
	err := db.transactionDB.Close()
	db.transactionDB = nil
	return err
}

// Rollback ...
func (db *Database) Rollback() error {
	if db.transactionDB == nil {
		return nil
	}

	return db.closeTransactionDB()
}

// HealthCheckWithErrorReturnWhenDown ...
func (db *Database) HealthCheckWithErrorReturnWhenDown() error {
	_, err := db.innerDB.Exec("SELECT 1")
	if err != nil {
		fmt.Println("PostgreSQL is down")
		return err
	}
	fmt.Println("PostgreSQL connection is ok")
	return nil
}
