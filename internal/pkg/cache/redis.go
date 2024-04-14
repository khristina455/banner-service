package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func (rc *RedisClient) Get(ctx context.Context, key string) ([]byte, bool) {
	value, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return value, true
}

func (rc *RedisClient) Set(key string, value []byte) {
	rc.client.Set(context.Background(), key, value, time.Minute*5)
}
