package cache

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return rdb
}
