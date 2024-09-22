package usecase

import (
	"bot_hmb/internal/db"
	"bot_hmb/internal/entity"
	"bot_hmb/internal/repository"
	"bot_hmb/internal/telegram"
	"context"
	"errors"
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

var (
	ErrIsNotMaster = errors.New("user not a master")
	ErrNotFound    = errors.New("not found")
)

type Repositories struct {
	UsersRepository           repository.UsersRepository
	UserSchoolsRepository     repository.UserSchoolsRepository
	UserPresentsRepository    repository.UserPresentsRepository
	SchoolsRepository         repository.SchoolsRepository
	SchoolsTrainingRepository repository.SchoolsTrainingRepository
	SubscriptionsRepository   repository.SubscriptionsRepository

	TelegramAccountsRepository    repository.TelegramAccountsRepository
	TelegramAuthTicketsRepository repository.TelegramAuthTicketsRepository

	StepRepository repository.StepRepository
}
type Usecase struct {
	Constructor telegram.Constructor
	Wrapper     telegram.Wrapper
	repo        Repositories

	config Config
}

type Config struct {
	TelegramBotURL     string
	MasterUserNickname string
}

func New(
	usersEncryptKey string,
	constructor telegram.Constructor,
	wrapper telegram.Wrapper,
	config Config,
	db *db.Manager,
	redis *db.RedisClient) *Usecase {
	return &Usecase{
		Wrapper:     wrapper,
		Constructor: constructor,
		config:      config,
		repo: Repositories{
			StepRepository:                repository.NewStepRepositoryWithRedis(redis),
			SubscriptionsRepository:       repository.NewSubscriptionsRepository(db.Gorm),
			SchoolsRepository:             repository.NewSchoolsRepository(db.Gorm),
			UserSchoolsRepository:         repository.NewUserSchoolsRepository(db.Gorm),
			UserPresentsRepository:        repository.NewUserPresentsRepository(db.Gorm),
			TelegramAuthTicketsRepository: repository.NewTelegramAuthTicketsRepository(db.Gorm),
			UsersRepository: repository.NewUsersRepositoryWithRedis(
				repository.NewUsersRepository(db.Gorm, usersEncryptKey),
				redis),
			SchoolsTrainingRepository: repository.NewSchoolsTrainingRepositoryWithRedis(
				repository.NewSchoolsTrainingRepository(db.Gorm),
				redis),
			TelegramAccountsRepository: repository.NewTelegramAccountsRepositoryWithRedis(
				repository.NewTelegramAccounts(db.Gorm), redis),
		},
	}
}

func (u *Usecase) getMaster(ctx context.Context, chatID int64) (entity.User, error) {
	masterUser, err := u.getUser(ctx, chatID)
	if err != nil {
		return entity.User{}, err
	}
	if masterUser.IsMaster || u.config.MasterUserNickname == masterUser.Username {
		return masterUser, nil
	}
	return entity.User{}, ErrIsNotMaster
}

func (u *Usecase) getUser(ctx context.Context, chatID int64) (entity.User, error) {
	var accountFound = false
	var user entity.User
	var userFound = false
	account, err := u.repo.TelegramAccountsRepository.FindActiveByChatID(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	} else if err == nil {
		accountFound = true
	}

	if accountFound {
		user, err = u.repo.UsersRepository.ByID(ctx, account.UserID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return user, err
		} else if err == nil {
			userFound = user.IsActivated && !user.IsDeleted
		}
	}
	if !userFound {
		return user, ErrNotFound
	}
	return user, nil
}

func checkPhone(phoneNumber string) (bool, error) {
	var pattern = `^7\d{10}$`
	re := regexp.MustCompile(pattern)
	res := re.FindAllString(phoneNumber, -1)
	if len(res) == 0 {
		return false, fmt.Errorf("phone number is not validation")
	}
	return true, nil
}
