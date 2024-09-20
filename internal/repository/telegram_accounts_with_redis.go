package repository

import (
	"bot_hmb/internal/db"
	"bot_hmb/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

const telegramAccountsCachePattern = "cache:telegram:accounts:%v"

type telegramAccountsRepositoryWithRedis struct {
	repo  TelegramAccountsRepository
	redis *db.RedisClient
}

func NewTelegramAccountsRepositoryWithRedis(repo TelegramAccountsRepository,
	redis *db.RedisClient) TelegramAccountsRepository {
	return &telegramAccountsRepositoryWithRedis{
		repo:  repo,
		redis: redis,
	}
}

func (r *telegramAccountsRepositoryWithRedis) FindActiveByChatID(
	ctx context.Context, chatID int64) (*entity.TelegramAccount, error) {
	var dto entity.TelegramAccount
	cacheData, err := r.redis.Get(ctx, fmt.Sprintf(telegramAccountsCachePattern, chatID))
	if err != nil {
		return &dto, err
	}
	if len(cacheData) == 0 {
		out, err := r.repo.FindActiveByChatID(ctx, chatID)
		if err != nil {
			return &dto, err
		}
		return out, r.setCacheTA(ctx, out)
	}
	err = json.Unmarshal(cacheData, &dto)
	if err != nil {
		return &dto, err
	}

	return &dto, nil
}

func (r *telegramAccountsRepositoryWithRedis) FindActiveByUserID(
	ctx context.Context, userID uuid.UUID) ([]*entity.TelegramAccount, error) {
	return r.repo.FindActiveByUserID(ctx, userID)
}

func (r *telegramAccountsRepositoryWithRedis) DetachChatID(ctx context.Context,
	chatID int64, userID uuid.UUID) error {
	err := r.redis.Del(ctx, fmt.Sprintf(telegramAccountsCachePattern, chatID))
	if err != nil {
		return err
	}
	return r.repo.DetachChatID(ctx, chatID, userID)
}

func (r *telegramAccountsRepositoryWithRedis) Upsert(ctx context.Context, account *entity.TelegramAccount) error {
	err := r.repo.Upsert(ctx, account)
	if err != nil {
		return err
	}
	return r.setCacheTA(ctx, account)
}

func (r *telegramAccountsRepositoryWithRedis) setCacheTA(ctx context.Context, dto *entity.TelegramAccount) error {
	var expiration = 2 * time.Minute
	cacheData, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, fmt.Sprintf(telegramAccountsCachePattern, dto.ChatID), cacheData, expiration)
	if err != nil {
		return err
	}
	return nil
}
