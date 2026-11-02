package cache

import (
	"context"
	"gateway/internal/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	Cache *redis.Client
}

func NewCacheClient(cfg *config.Config) *CacheClient {
	options := &redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,

		MinIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}

	rdb := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	log.Println("redis connect success")
	return &CacheClient{
		Cache: rdb,
	}
}
