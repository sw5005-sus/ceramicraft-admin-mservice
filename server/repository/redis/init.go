package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/config"
)

var redisClient *redis.Client

func InitRedis() {
	address := fmt.Sprintf("%s:%d", config.Config.RedisConfig.Host, config.Config.RedisConfig.Port)
	redisClient = redis.NewClient(&redis.Options{
		Addr: address,
		DB:   0, // Use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
