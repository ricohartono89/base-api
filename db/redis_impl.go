package db

import (
	"time"
)

// GetRedisValue ...
func (r *Redis) GetRedisValue(redisKey string) (string, error) {
	redisValue, err := r.client.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	return redisValue, nil
}

// SetRedisValue ...
func (r *Redis) SetRedisValue(redisKey string, redisValue interface{}, timeOutPeriod time.Duration) {
	r.client.Set(redisKey, redisValue, timeOutPeriod)
}

// DeleteRedisValue ...
func (r *Redis) DeleteRedisValue(redisKey string) {
	r.client.Del(redisKey)
}
