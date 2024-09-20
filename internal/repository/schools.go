package repository

import (
	"bot_hmb/internal/entity"
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	schoolsTable = "schools"
)

type SchoolsRepository interface {
	Create(ctx context.Context, user *entity.School) error
	ByID(ctx context.Context, id uuid.UUID) (entity.School, error)
	ByRegion(ctx context.Context, region string) (entity.School, error)
	List(ctx context.Context) ([]entity.School, error)
}

type schoolsRepository struct {
	Db *gorm.DB
}

func NewSchoolsRepository(db *gorm.DB) SchoolsRepository {
	return &schoolsRepository{
		Db: db,
	}
}

func (r *schoolsRepository) Create(ctx context.Context, school *entity.School) error {
	result := r.Db.WithContext(ctx).
		Table(schoolsTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(school)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *schoolsRepository) ByID(ctx context.Context, id uuid.UUID) (entity.School, error) {
	dto := entity.School{}
	result := r.Db.WithContext(ctx).
		Table(schoolsTable).
		Where("id = ?", id).
		Take(&dto)

	if result.Error != nil {
		return dto, result.Error
	}

	return dto, nil
}

func (r *schoolsRepository) ByRegion(ctx context.Context, region string) (entity.School, error) {
	dto := entity.School{}
	result := r.Db.WithContext(ctx).
		Table(schoolsTable).
		Where("region = ?", region).
		Take(&dto)

	if result.Error != nil {
		return dto, result.Error
	}

	return dto, nil
}

func (r *schoolsRepository) List(ctx context.Context) ([]entity.School, error) {
	dtos := make([]entity.School, 0)
	result := r.Db.WithContext(ctx).
		Table(schoolsTable).
		Find(&dtos)

	if result.Error != nil {
		return dtos, result.Error
	}

	return dtos, nil
}
