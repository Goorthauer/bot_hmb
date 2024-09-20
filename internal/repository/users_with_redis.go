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

const userCachePattern = "cache:user:%s"

type usersRepositoryWithRedis struct {
	repo  UsersRepository
	redis *db.RedisClient
}

func NewUsersRepositoryWithRedis(repo UsersRepository, redis *db.RedisClient) UsersRepository {
	return &usersRepositoryWithRedis{
		repo:  repo,
		redis: redis,
	}
}

func (r *usersRepositoryWithRedis) Create(ctx context.Context, user *entity.User) error {
	return r.repo.Create(ctx, user)
}

func (r *usersRepositoryWithRedis) ByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	user, err := r.getUserFromCacheOrDB(ctx, fmt.Sprintf(userCachePattern, userID), func() (entity.User, error) {
		return r.repo.ByID(ctx, userID)
	})
	return user, err
}

func (r *usersRepositoryWithRedis) BySchool(ctx context.Context, schoolID uuid.UUID) ([]*entity.User, error) {
	return r.repo.BySchool(ctx, schoolID)
}

func (r *usersRepositoryWithRedis) ByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entity.User, error) {
	users, missingUserIDs, err := r.getUsersFromCache(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	if len(missingUserIDs) > 0 {
		missingUsers, err := r.repo.ByIDs(ctx, missingUserIDs)
		if err != nil {
			return nil, err
		}
		for _, user := range missingUsers {
			err = r.setCacheUser(ctx, user)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
	}

	return users, nil
}

func (r *usersRepositoryWithRedis) ByUsername(ctx context.Context, username string) (entity.User, error) {
	user, err := r.getUserFromCacheOrDB(ctx, fmt.Sprintf(userCachePattern, username), func() (entity.User, error) {
		return r.repo.ByUsername(ctx, username)
	})
	return user, err
}

func (r *usersRepositoryWithRedis) ByPhone(ctx context.Context, phone string) (entity.User, error) {
	user, err := r.getUserFromCacheOrDB(ctx, fmt.Sprintf(userCachePattern, phone), func() (entity.User, error) {
		return r.repo.ByPhone(ctx, phone)
	})
	return user, err
}

func (r *usersRepositoryWithRedis) SetMasterRights(ctx context.Context, userID uuid.UUID, isMaster bool) error {
	err := r.redis.Del(ctx, fmt.Sprintf(userCachePattern, userID))
	if err != nil {
		return err
	}
	return r.repo.SetMasterRights(ctx, userID, isMaster)
}

func (r *usersRepositoryWithRedis) getUserFromCacheOrDB(ctx context.Context, cacheKey string, fetchFromDB func() (entity.User, error)) (entity.User, error) {
	var user entity.User
	cacheData, err := r.redis.Get(ctx, cacheKey)
	if err != nil {
		return user, err
	}

	if len(cacheData) > 0 {
		err = json.Unmarshal(cacheData, &user)
		if err != nil {
			return user, err
		}
		return user, nil
	}

	user, err = fetchFromDB()
	if err != nil {
		return user, err
	}

	err = r.setCacheUser(ctx, &user)
	return user, err
}

func (r *usersRepositoryWithRedis) getUsersFromCache(ctx context.Context, userIDs []uuid.UUID) ([]*entity.User, []uuid.UUID, error) {
	users := make([]*entity.User, 0, len(userIDs))
	missingUserIDs := make([]uuid.UUID, 0)

	for _, userID := range userIDs {
		cacheData, err := r.redis.Get(ctx, fmt.Sprintf(userCachePattern, userID))
		if err != nil || len(cacheData) == 0 {
			missingUserIDs = append(missingUserIDs, userID)
			continue
		}

		var user entity.User
		err = json.Unmarshal(cacheData, &user)
		if err != nil {
			return nil, nil, err
		}
		users = append(users, &user)
	}

	return users, missingUserIDs, nil
}

func (r *usersRepositoryWithRedis) setCacheUser(ctx context.Context, user *entity.User) error {
	cacheData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	expiration := 2 * time.Minute
	err = r.redis.Set(ctx, fmt.Sprintf(userCachePattern, user.ID), cacheData, expiration)
	if err != nil {
		return err
	}

	if user.Username != "" {
		err = r.redis.Set(ctx, fmt.Sprintf(userCachePattern, user.Username), cacheData, expiration)
		if err != nil {
			return err
		}
	}

	if user.Phone != "" {
		err = r.redis.Set(ctx, fmt.Sprintf(userCachePattern, user.Phone), cacheData, expiration)
		if err != nil {
			return err
		}
	}

	return nil
}
