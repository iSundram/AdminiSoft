
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
package config

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis(cfg *Config) *redis.Client {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		// Fallback to default settings
		opt = &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}
	}

	rdb := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		// Handle connection error gracefully
		return nil
	}

	return rdb
}
