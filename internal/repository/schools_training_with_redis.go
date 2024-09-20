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

const schoolsTrainingPattern = "cache:school:training:%v"

type schoolsTrainingRepositoryWithRedis struct {
	repo  SchoolsTrainingRepository
	redis *db.RedisClient
}

func NewSchoolsTrainingRepositoryWithRedis(repo SchoolsTrainingRepository,
	redis *db.RedisClient) SchoolsTrainingRepository {
	return &schoolsTrainingRepositoryWithRedis{
		repo:  repo,
		redis: redis,
	}
}

func (r *schoolsTrainingRepositoryWithRedis) Create(ctx context.Context, dto *entity.SchoolTraining) error {
	err := r.redis.Del(ctx, fmt.Sprintf(schoolsTrainingPattern, dto.SchoolID))
	if err != nil {
		return err
	}
	return r.repo.Create(ctx, dto)
}
func (r *schoolsTrainingRepositoryWithRedis) BySchool(ctx context.Context, schoolID uuid.UUID) (entity.SchoolTraining, error) {
	dto := entity.SchoolTraining{}
	cacheData, err := r.redis.Get(ctx, fmt.Sprintf(schoolsTrainingPattern, schoolID))
	if err != nil {
		return dto, err
	}
	if len(cacheData) == 0 {
		dto, err = r.repo.BySchool(ctx, schoolID)
		return dto, r.setCache(ctx, &dto)
	}
	err = json.Unmarshal(cacheData, &dto)
	if err != nil {
		return dto, err
	}

	return dto, nil
}

func (r *schoolsTrainingRepositoryWithRedis) setCache(ctx context.Context, dto *entity.SchoolTraining) error {
	var expiration = 8 * time.Hour
	cacheData, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, fmt.Sprintf(userCachePattern, dto.SchoolID), cacheData, expiration)
	if err != nil {
		return err
	}
	return nil
}
