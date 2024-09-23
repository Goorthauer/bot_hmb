package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func (u *Usecase) Registration(ctx context.Context, chatID int64, username string) error {
	account, err := u.repo.TelegramAccountsRepository.FindActiveByChatID(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if account.ChatID != 0 {
		user, err := u.repo.UsersRepository.ByID(ctx, account.UserID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return u.Constructor.
			ConstructStartAttachedAndSend(ctx,
				[]int64{chatID},
				user.GetFullNameWithSchool())
	}
	if username != "" {
		user, err := u.repo.UsersRepository.ByUsername(ctx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return u.Constructor.
				ConstructStartAttachedAndSend(ctx,
					[]int64{chatID},
					user.GetFullNameWithSchool())
		}
	}

	step := &entity.Step{
		ChatID:   chatID,
		Step:     entity.StepPhone,
		Username: username,
	}
	err = u.repo.StepRepository.SetStep(ctx, step)
	if err != nil {
		return err
	}

	return u.Constructor.
		ConstructRegisterWithStep(ctx,
			chatID,
			step.Step,
			nil)
}

func (u *Usecase) RegistrationByStep(ctx context.Context, chatID int64, field string) error {
	if field == "" {
		return errors.New("data is empty")
	}
	field = strings.TrimSpace(field)
	step, err := u.repo.StepRepository.ByChatID(ctx, chatID)
	if err != nil {
		return nil
	}
	if !step.Step.IsValid() {
		return nil
	}
	schoolList := make(map[uuid.UUID]string)
	switch step.Step {
	case entity.StepFirstName:
		if _, err := checkPhone(step.Phone); err != nil {
			_ = u.repo.StepRepository.DelStep(ctx, chatID)
			return err
		}
		step.Firstname = field
	case entity.StepLastname:
		step.Lastname = field
		schools, err := u.repo.SchoolsRepository.List(ctx)
		if err != nil {
			return err
		}
		for _, v := range schools {
			schoolList[v.ID] = v.Name
		}
	case entity.StepPhone:
		step.Phone = field
		user, err := u.repo.UsersRepository.ByPhone(ctx, field)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if user.IsActivated && !user.IsDeleted {
			err = u.repo.StepRepository.DelStep(ctx, chatID)
			if err != nil {
				return err
			}
			return u.ConnectTelegramToUserID(ctx, chatID, user.ID, user.SchoolName)
		}
	case entity.StepSchool:
		return u.repo.StepRepository.DelStep(ctx, chatID)
	default:
		return nil

	}

	step.Step = step.Step.Next()

	err = u.repo.StepRepository.SetStep(ctx, &step)
	if err != nil {
		return err
	}
	return u.Constructor.
		ConstructRegisterWithStep(ctx,
			chatID,
			step.Step,
			schoolList)
}

func (u *Usecase) RegistrationLastStep(ctx context.Context, chatID int64,
	schoolID string) error {
	step, err := u.repo.StepRepository.ByChatID(ctx, chatID)
	user, err := u.registerUser(ctx, schoolID, step.Firstname, step.Lastname, step.Phone, step.Username)
	if err != nil {
		return err
	}
	err = u.ConnectTelegramToUserID(ctx, chatID, user.ID, user.SchoolName)
	if err != nil {
		return err
	}
	err = u.repo.StepRepository.DelStep(ctx, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) registerUser(ctx context.Context,
	schoolID, firstname, lastname, phoneNumber, username string) (*entity.User, error) {
	user, err := u.repo.UsersRepository.ByPhone(ctx, phoneNumber)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

	}
	if user.ID == uuid.Nil {
		now := time.Now()
		userID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		randomBytes := make([]byte, 256)
		_, err = rand.Read(randomBytes)
		if err != nil {
			return nil, err
		}
		pdEncryptionKey := base64.StdEncoding.EncodeToString(randomBytes)
		user = entity.User{
			ID:              userID,
			IsActivated:     true,
			RegisteredAt:    now,
			Username:        username,
			Phone:           phoneNumber,
			PDEncryptionKey: &pdEncryptionKey,
			PersonalData: entity.UserPersonalData{
				Firstname: &firstname,
				Lastname:  &lastname,
			},
		}
		err = u.repo.UsersRepository.Create(ctx, &user)
		if err != nil {
			return nil, err
		}
	}
	school, err := u.repo.SchoolsRepository.ByID(ctx, uuid.FromStringOrNil(schoolID))
	if err != nil {
		return nil, err
	}

	user.SchoolName = school.Name
	return &user, u.repo.UserSchoolsRepository.Create(ctx, &entity.UserSchool{
		UserID:   user.ID,
		SchoolID: school.ID,
	})
}
