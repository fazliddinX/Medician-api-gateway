package redis

import (
	config2 "api-gateway/pkg/config"
	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"
)

func RabbitClient(cfg config2.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RADIS + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}
