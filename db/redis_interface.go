package db

import "time"

// RedisInterface ...
type RedisInterface interface {
	Connect()
	Close()
	GetRedisValue(redisKey string) (string, error)
	SetRedisValue(redisKey string, redisValue interface{}, timeOutDuration time.Duration)
	DeleteRedisValue(redisKey string)
}
