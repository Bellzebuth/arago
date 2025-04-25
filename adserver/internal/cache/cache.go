package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
)

func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return rdb
}
