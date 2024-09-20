package db

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, addr string) (*RedisClient, error) {
	r := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   2,
	})
	if err := r.Set(ctx, "key", "value", 0).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{client: r}, nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if err := rc.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}
func (rc *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	out, err := rc.client.Get(ctx, key).Bytes()
	if err != redis.Nil {
		return out, err
	}
	return out, nil
}
func (rc *RedisClient) Del(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}

func (rc *RedisClient) Incr(ctx context.Context, key string) (int, error) {
	count, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		return int(count), err
	}
	return int(count), nil
}

func (rc *RedisClient) Expire(ctx context.Context, key string, interval time.Duration) error {
	return rc.client.Expire(ctx, key, interval).Err()
}
