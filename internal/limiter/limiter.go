package limiter

import (
	"bot_hmb/internal/db"
	"context"
	"fmt"
	"time"
)

type ChatRateLimiter interface {
	AllowRequest(ctx context.Context, chatID int64) (bool, error)
}

type chatRateLimiter struct {
	redisClient *db.RedisClient
	limit       int
	interval    time.Duration
}

type RateLimit struct {
	lastRequest time.Time
	count       int
}

func NewRateLimiter(redisClient *db.RedisClient, limitPerChat int, interval time.Duration) ChatRateLimiter {
	return &chatRateLimiter{
		redisClient: redisClient,
		limit:       limitPerChat,
		interval:    interval,
	}
}

func (r *chatRateLimiter) AllowRequest(ctx context.Context, chatID int64) (bool, error) {
	const cachePattern = "chat:%d:requests"
	key := fmt.Sprintf(cachePattern, chatID)

	count, err := r.redisClient.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = r.redisClient.Expire(ctx, key, r.interval)
		if err != nil {
			return false, err
		}
	}
	if count > r.limit {
		return false, nil
	}

	return true, nil
}
