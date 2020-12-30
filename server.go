package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ricohartono89/base-api/db"
	"github.com/ricohartono89/base-api/env"
	"github.com/ricohartono89/base-api/internal/router/way"

	"github.com/ricohartono89/base-api/middleware"
)

type server struct {
	router *way.Router
	db     *db.DatabaseInterface
	redis  *db.RedisInterface
	client *http.Client
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	s.router.ServeHTTP(w, r)
}

func newServer(router *way.Router, dbi db.DatabaseInterface, redis db.RedisInterface) *server {
	httpClient := &http.Client{Timeout: time.Second * 10}
	s := &server{
		router: router,
		db:     &dbi,
		redis:  &redis,
	}

	s.routes(httpClient, middleware.Middleware{Client: httpClient})

	return s
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// flags
	addr := flag.String("addr", env.AppPort(), "http service address")
	flag.Parse()

	database := db.Database{}
	database.Connect()
	defer database.Close()

	client := db.Redis{}
	client.Connect()
	defer client.Close()

	srv := newServer(way.NewRouter(), &database, &client)
	s := &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         *addr,
		Handler:      srv.router,
	}

	return s.ListenAndServe()
}
