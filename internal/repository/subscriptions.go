package repository

import (
	"bot_hmb/internal/entity"
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	subscriptionTable = "subscriptions"
)

type SubscriptionsRepository interface {
	Create(ctx context.Context, dto *entity.Subscription) error
	ByUserID(ctx context.Context, userID uuid.UUID) (entity.Subscription, error)
	BySchoolID(ctx context.Context, schoolID uuid.UUID, onlyLast bool) ([]entity.Subscription, error)
}

type subscriptionsRepository struct {
	Db *gorm.DB
}

func NewSubscriptionsRepository(db *gorm.DB) SubscriptionsRepository {
	return &subscriptionsRepository{
		Db: db,
	}
}

func (r *subscriptionsRepository) Create(ctx context.Context, dto *entity.Subscription) error {
	result := r.Db.WithContext(ctx).
		Table(subscriptionTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *subscriptionsRepository) ByUserID(ctx context.Context, userID uuid.UUID) (entity.Subscription, error) {
	dto := entity.Subscription{}
	result := r.Db.WithContext(ctx).
		Table(subscriptionTable).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Take(&dto)

	if result.Error != nil {
		return dto, result.Error
	}

	return dto, nil
}

func (r *subscriptionsRepository) BySchoolID(ctx context.Context, schoolID uuid.UUID, onlyLast bool) ([]entity.Subscription, error) {
	out := make([]entity.Subscription, 0)

	result := r.Db.WithContext(ctx).
		Table(subscriptionTable).
		Where("school_id = ?", schoolID).
		Order("deadline_at")

	if onlyLast {
		subQuery := r.Db.WithContext(ctx).Table(subscriptionTable).
			Select(`DISTINCT ON (user_id) id`).
			Order(`user_id, created_at DESC`)
		result.Where("id IN (?)", subQuery)
	}
	result = result.Find(&out)

	if result.Error != nil {
		return out, result.Error
	}

	return out, nil
}
