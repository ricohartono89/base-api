package db

import (
	"crypto/tls"
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/ricohartono89/base-api/env"
)

// Redis ...
type Redis struct {
	client *redis.Client
}

// Connect ...
func (r *Redis) Connect() {
	var tlsConfig *tls.Config
	if env.Env() == "production" || env.Env() == "staging" || env.Env() == "development" {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", env.RedisUrl(), env.RedisPort()),
		Password:  env.RedisPassword(),
		TLSConfig: tlsConfig,
	})

	_, err := r.client.Ping().Result()
	if err != nil {
		fmt.Println("Failed to establish redis connection!")
		panic(err)
	}
}

// Close ...
func (r *Redis) Close() {
	r.client.Close()
}
