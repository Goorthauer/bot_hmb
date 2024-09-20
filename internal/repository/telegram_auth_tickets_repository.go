package repository

import (
	entity "bot_hmb/internal/entity"
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	telegramAuthTicketsTable = "telegram_auth_tickets"
)

type TelegramAuthTicketsRepository interface {
	ByToken(ctx context.Context, token string) (*entity.TelegramAuthTicket, error)
	Upsert(ctx context.Context, ticket *entity.TelegramAuthTicket) error
	FindActiveTickets(ctx context.Context, userID uuid.UUID) ([]*entity.TelegramAuthTicket, error)
	DeactivateTicket(ctx context.Context, token string) error
	SpendTicket(ctx context.Context, token string, now time.Time) error
}

type telegramAuthTicketsRepository struct {
	Db *gorm.DB
}

func NewTelegramAuthTicketsRepository(db *gorm.DB) TelegramAuthTicketsRepository {
	return &telegramAuthTicketsRepository{
		Db: db,
	}
}

func (r *telegramAuthTicketsRepository) ByToken(
	ctx context.Context, token string) (*entity.TelegramAuthTicket, error) {
	ticket := &entity.TelegramAuthTicket{}
	result := r.Db.WithContext(ctx).
		Table(telegramAuthTicketsTable).
		Where("token = ?", token).
		Where("expires_at > now()").
		Not("is_blocked").
		Not("is_spent").
		Take(ticket)

	if result.Error != nil {
		return ticket, result.Error
	}

	return ticket, nil
}

func (r *telegramAuthTicketsRepository) Upsert(ctx context.Context, ticket *entity.TelegramAuthTicket) error {
	result := r.Db.WithContext(ctx).
		Table(telegramAuthTicketsTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "token"}},
			UpdateAll: true,
		}).
		Create(ticket)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *telegramAuthTicketsRepository) FindActiveTickets(
	ctx context.Context, userID uuid.UUID) ([]*entity.TelegramAuthTicket, error) {
	tickets := make([]*entity.TelegramAuthTicket, 0)

	result := r.Db.WithContext(ctx).
		Table(telegramAuthTicketsTable).
		Joins(fmt.Sprintf("JOIN %[1]s ON %[1]s.id = %[2]s.user_id",
			usersTable,
			telegramAuthTicketsTable)).
		Where(fmt.Sprintf("%s.id = ?", usersTable), userID).
		Not(fmt.Sprintf("%s.is_blocked", telegramAuthTicketsTable)).
		Not(fmt.Sprintf("%s.is_spent", telegramAuthTicketsTable)).
		Find(&tickets)

	if result.Error != nil {
		return nil, result.Error
	}

	return tickets, nil
}

func (r *telegramAuthTicketsRepository) DeactivateTicket(ctx context.Context, token string) error {
	result := r.Db.WithContext(ctx).
		Table(telegramAuthTicketsTable).
		Where("token = ?", token).
		Not("is_blocked").
		UpdateColumn("is_blocked", true)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *telegramAuthTicketsRepository) SpendTicket(ctx context.Context, token string, now time.Time) error {
	result := r.Db.WithContext(ctx).
		Table(telegramAuthTicketsTable).
		Where("token = ?", token).
		Not("is_spent").
		Updates(map[string]interface{}{
			"is_spent": true,
			"spent_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
