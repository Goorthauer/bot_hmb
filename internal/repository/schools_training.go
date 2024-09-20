package repository

import (
	"bot_hmb/internal/entity"
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	schoolsTrainingTable = "schools_training"
)

type SchoolsTrainingRepository interface {
	Create(ctx context.Context, dto *entity.SchoolTraining) error
	BySchool(ctx context.Context, schoolID uuid.UUID) (entity.SchoolTraining, error)
}

type schoolsTrainingRepository struct {
	Db *gorm.DB
}

func NewSchoolsTrainingRepository(db *gorm.DB) SchoolsTrainingRepository {
	return &schoolsTrainingRepository{
		Db: db,
	}
}

func (r *schoolsTrainingRepository) Create(ctx context.Context, dto *entity.SchoolTraining) error {
	result := r.Db.WithContext(ctx).
		Table(schoolsTrainingTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "school_id"}, {Name: "created_at"}},
			UpdateAll: true,
		}).
		Create(dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (r *schoolsTrainingRepository) BySchool(ctx context.Context, schoolID uuid.UUID) (entity.SchoolTraining, error) {
	user := entity.SchoolTraining{}
	result := r.Db.WithContext(ctx).
		Table(schoolsTrainingTable).
		Where("school_id = ?", schoolID).
		Order("created_at desc").
		Take(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
