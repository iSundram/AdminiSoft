
package config

import (
	"github.com/go-redis/redis/v8"
	"context"
)

func InitRedis(cfg *Config) *redis.Client {
	opts, _ := redis.ParseURL(cfg.RedisURL)
	client := redis.NewClient(opts)
	
	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	
	return client
}
