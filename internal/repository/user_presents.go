package repository

import (
	"bot_hmb/internal/entity"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	userPresentsTable = "user_presents"
)

type UserPresentsRepository interface {
	Create(ctx context.Context, dto *entity.Presents) error
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
		Debug().
		Create(dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
