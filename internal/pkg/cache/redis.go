package cache

import (
	"context"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

// TODO: поработать с неймингом

type RedisClient struct {
	client      *cache.Cache
	redisClient *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})

	return &RedisClient{client: cache, redisClient: client}
}

func (rc *RedisClient) Get(ctx context.Context, key string) (value []byte, ok bool) {
	value, err := rc.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return value, true
}

func (rc *RedisClient) Set(key string, value []byte) {
	rc.redisClient.Set(context.Background(), key, value, time.Minute*5)
}
