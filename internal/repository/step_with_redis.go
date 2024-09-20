package repository

import (
	"bot_hmb/internal/db"
	"bot_hmb/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const stepCachePattern = "cache:invite:step:%v"

type StepRepository interface {
	DelStep(ctx context.Context, chatID int64) error
	ByChatID(ctx context.Context, chatID int64) (entity.Step, error)
	SetStep(ctx context.Context, dto *entity.Step) error
}
type stepRepositoryWithRedis struct {
	redis *db.RedisClient
}

func NewStepRepositoryWithRedis(
	redis *db.RedisClient) StepRepository {
	return &stepRepositoryWithRedis{
		redis: redis,
	}
}

func (r *stepRepositoryWithRedis) ByChatID(ctx context.Context, chatID int64) (entity.Step, error) {
	dto := entity.Step{}
	cacheData, err := r.redis.Get(ctx, fmt.Sprintf(stepCachePattern, chatID))
	if err != nil {
		return dto, err
	}
	if len(cacheData) != 0 {
		err = json.Unmarshal(cacheData, &dto)
		if err != nil {
			return dto, err
		}
	}
	return dto, nil
}

func (r *stepRepositoryWithRedis) SetStep(ctx context.Context, dto *entity.Step) error {
	var expiration = 1 * time.Hour
	cacheData, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, fmt.Sprintf(stepCachePattern, dto.ChatID), cacheData, expiration)
	if err != nil {
		return err
	}
	return nil
}

func (r *stepRepositoryWithRedis) DelStep(ctx context.Context, chatID int64) error {
	return r.redis.Del(ctx, fmt.Sprintf(stepCachePattern, chatID))
}
