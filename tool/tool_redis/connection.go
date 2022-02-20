package tool_redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var GlobClient *redis.Client

func NewRedis() error {
	var opt = redis.Options{
		Addr: "101.35.193.209:6380",
	}
	GlobClient = redis.NewClient(&opt)
	if err := GlobClient.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}
