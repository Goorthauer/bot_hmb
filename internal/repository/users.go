package repository

import (
	"bot_hmb/internal/entity"
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const usersTable = "users"

type UsersRepository interface {
	Create(ctx context.Context, user *entity.User) error
	ByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	ByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entity.User, error)
	BySchool(ctx context.Context, schoolID uuid.UUID, onlySubscription bool) ([]*entity.User, error)
	ByUsername(ctx context.Context, username string) (entity.User, error)
	ByPhone(ctx context.Context, phone string) (entity.User, error)
	SetMasterRights(ctx context.Context, userID uuid.UUID, isMaster bool) error
}

type usersRepository struct {
	Db                       *gorm.DB
	usersPDBaseEncryptionKey string
}

func NewUsersRepository(db *gorm.DB, usersPDBaseEncryptionKey string) UsersRepository {
	return &usersRepository{
		Db:                       db,
		usersPDBaseEncryptionKey: usersPDBaseEncryptionKey,
	}
}

func (r *usersRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.encryptPersonalData(user); err != nil {
		return err
	}

	result := r.Db.WithContext(ctx).
		Table(usersTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(user)
	return result.Error
}

func (r *usersRepository) ByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	return r.fetchUser(ctx, fmt.Sprintf("%s.id = ?", usersTable), userID)
}

func (r *usersRepository) BySchool(ctx context.Context,
	schoolID uuid.UUID,
	onlySubscription bool) ([]*entity.User, error) {
	var users []*entity.User
	tx := r.baseTx(ctx)
	if onlySubscription {
		subQuery := r.getSubscriptionTx(ctx)
		tx.Joins(
			fmt.Sprintf("JOIN (?) as %[1]s on %[2]s.id = %[1]s.user_id", subscriptionTable, usersTable),
			subQuery)
	}
	result := tx.Where(fmt.Sprintf("%s.school_id = ?", userSchoolsTable), schoolID).
		Where("is_activated").
		Not("is_deleted").
		Debug().
		Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	for _, user := range users {
		if err := r.decryptPersonalData(user); err != nil {
			return users, err
		}
	}
	return users, nil
}

func (r *usersRepository) ByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entity.User, error) {
	var users []*entity.User
	result := r.baseTx(ctx).
		Where("users.id = ANY(?)", pq.Array(userIDs)).
		Not("is_deleted").
		Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	for _, user := range users {
		if err := r.decryptPersonalData(user); err != nil {
			return users, err
		}
	}
	return users, nil
}

func (r *usersRepository) ByUsername(ctx context.Context, username string) (entity.User, error) {
	return r.fetchUser(ctx, fmt.Sprintf("%s.username = ?", usersTable), username)
}

func (r *usersRepository) ByPhone(ctx context.Context, phone string) (entity.User, error) {
	return r.fetchUser(ctx, fmt.Sprintf("%s.phone = ?", usersTable), phone)
}

func (r *usersRepository) SetMasterRights(ctx context.Context, userID uuid.UUID, isMaster bool) error {
	result := r.Db.WithContext(ctx).
		Table(usersTable).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_master": isMaster,
		})
	return result.Error
}

func (r *usersRepository) baseTx(ctx context.Context) *gorm.DB {
	return r.Db.WithContext(ctx).
		Select(fmt.Sprintf(
			"%[1]s.*,%[2]s.region,%[2]s.id as school_id,%[2]s.name as school_name,%[2]s.address as school_address",
			usersTable, schoolsTable)).
		Table(usersTable).
		Joins(fmt.Sprintf("JOIN %[1]s ON %[1]s.user_id = %[2]s.id", userSchoolsTable, usersTable)).
		Joins(fmt.Sprintf("JOIN %[2]s ON %[1]s.school_id = %[2]s.id", userSchoolsTable, schoolsTable))
}

func (r *usersRepository) fetchUser(ctx context.Context, query interface{}, args ...interface{}) (entity.User, error) {
	user := entity.User{}
	result := r.baseTx(ctx).Where(query, args...).Not("is_deleted").Take(&user)
	if result.Error != nil {
		return user, result.Error
	}
	if err := r.decryptPersonalData(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (r *usersRepository) encryptPersonalData(user *entity.User) error {
	return user.PersonalData.Encrypt(r.usersPDBaseEncryptionKey, user.PDEncryptionKey)
}

func (r *usersRepository) decryptPersonalData(user *entity.User) error {
	return user.PersonalData.Decrypt(r.usersPDBaseEncryptionKey, user.PDEncryptionKey)
}

func (r *usersRepository) getSubscriptionTx(ctx context.Context) *gorm.DB {
	txSub := r.Db.WithContext(ctx).Table(subscriptionTable).
		Select(`DISTINCT ON (user_id) id`).
		Order(`user_id, created_at DESC`)
	subQuery := r.Db.WithContext(ctx).Table(subscriptionTable).Where("id IN (?)", txSub)
	return subQuery
}
