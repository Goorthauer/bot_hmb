package repository

import (
	"bot_hmb/internal/entity"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	userSchoolsTable = "user_schools"
)

type UserSchoolsRepository interface {
	Create(ctx context.Context, dto *entity.UserSchool) error
}

type userSchoolsRepository struct {
	Db *gorm.DB
}

func NewUserSchoolsRepository(db *gorm.DB) UserSchoolsRepository {
	return &userSchoolsRepository{
		Db: db,
	}
}

func (r *userSchoolsRepository) Create(ctx context.Context, dto *entity.UserSchool) error {
	result := r.Db.WithContext(ctx).
		Table(userSchoolsTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).
		Debug().
		Create(dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
