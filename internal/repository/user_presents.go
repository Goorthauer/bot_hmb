package repository

import (
	"bot_hmb/internal/entity"
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	userPresentsTable = "user_presents"
)

type UserPresentsRepository interface {
	Create(ctx context.Context, dto *entity.Presents) error
	ByUserID(ctx context.Context, userID uuid.UUID, lastCount int) ([]*entity.Presents, error)

	BySchoolIDWithLost(ctx context.Context, schoolID uuid.UUID) ([]*entity.PresentsSubscription, error)
}

type userPresentsRepository struct {
	Db *gorm.DB
}

func NewUserPresentsRepository(db *gorm.DB) UserPresentsRepository {
	return &userPresentsRepository{
		Db: db,
	}
}

func (r *userPresentsRepository) Create(ctx context.Context, dto *entity.Presents) error {
	result := r.Db.WithContext(ctx).
		Table(userPresentsTable).
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

func (r *userPresentsRepository) ByUserID(ctx context.Context, userID uuid.UUID, lastCount int) ([]*entity.Presents, error) {
	var dto []*entity.Presents
	result := r.Db.
		WithContext(ctx).
		Table(userPresentsTable).
		Where(`user_id = ?`, userID).
		Limit(lastCount).
		Find(&dto)
	if result.Error != nil {
		return dto, result.Error
	}
	return dto, nil
}

func (r *userPresentsRepository) BySchoolIDWithLost(ctx context.Context, schoolID uuid.UUID) ([]*entity.PresentsSubscription, error) {
	var dto []*entity.PresentsSubscription
	subTx := r.Db.
		WithContext(ctx).
		Select(`user_id,sum(days) as days,min(created_at) as created_at,max(deadline_at) as deadline_at`).
		Table(subscriptionTable).
		Where(`now() between (date_trunc('day',created_at)) and deadline_at`).
		Where(`school_id = ?`, schoolID).
		Group(`user_id`)
	result := r.Db.WithContext(ctx).
		Select(fmt.Sprintf(
			`%[1]s.user_id,
					count(%[1]s.id) as count_training,
					%[2]s.days as subscription_days,
					%[2]s.deadline_at as deadline_at`,
			userPresentsTable, subscriptionTable)).
		Table(userPresentsTable).
		Joins(fmt.Sprintf(`JOIN (?) %[2]s on %[2]s.user_id = %[1]s.user_id
and %[1]s.created_at between date_trunc('day', %[2]s.created_at) and %[2]s.deadline_at`,
			userPresentsTable, subscriptionTable), subTx).
		Group(fmt.Sprintf(`%[1]s.user_id, %[2]s.days, %[2]s.deadline_at`, userPresentsTable, subscriptionTable)).
		Find(&dto)
	if result.Error != nil {
		return dto, result.Error
	}
	return dto, nil
}
