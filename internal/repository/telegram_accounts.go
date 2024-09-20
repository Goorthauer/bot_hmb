package repository

import (
	"bot_hmb/internal/entity"
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	telegramAccountsTable = "telegram_accounts"
)

type TelegramAccountsRepository interface {
	FindActiveByChatID(ctx context.Context, chatID int64) (*entity.TelegramAccount, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.TelegramAccount, error)
	DetachChatID(ctx context.Context, chatID int64, userID uuid.UUID) error
	Upsert(ctx context.Context, account *entity.TelegramAccount) error
}

func NewTelegramAccounts(db *gorm.DB) TelegramAccountsRepository {
	return &telegramAccountsRepository{Db: db}
}

type telegramAccountsRepository struct {
	Db *gorm.DB
}

func (r *telegramAccountsRepository) FindActiveByChatID(
	ctx context.Context, chatID int64) (*entity.TelegramAccount, error) {
	account := &entity.TelegramAccount{}
	result := r.Db.WithContext(ctx).
		Table(telegramAccountsTable).
		Where("chat_id = ?", chatID).
		Where("is_active").
		Take(account)

	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func (r *telegramAccountsRepository) FindActiveByUserID(
	ctx context.Context, userID uuid.UUID) ([]*entity.TelegramAccount, error) {
	accounts := make([]*entity.TelegramAccount, 0)
	result := r.Db.WithContext(ctx).
		Table(telegramAccountsTable).
		Where("user_id = ?", userID).
		Where("is_active").
		Order("created_at DESC").
		Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (r *telegramAccountsRepository) DetachChatID(ctx context.Context,
	chatID int64, userID uuid.UUID) error {
	result := r.Db.WithContext(ctx).
		Table(telegramAccountsTable).
		Where("user_id = ?", userID).
		Where("chat_id = ?", chatID).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *telegramAccountsRepository) Upsert(ctx context.Context, account *entity.TelegramAccount) error {
	result := r.Db.WithContext(ctx).
		Table(telegramAccountsTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "chat_id"}},
			UpdateAll: true,
		}).
		Create(account)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
