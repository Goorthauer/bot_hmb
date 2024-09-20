package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (u *Usecase) StartCommand(ctx context.Context,
	chatID int64,
	ticket,
	username string) error {
	now := time.Now()

	account, err := u.repo.TelegramAccountsRepository.FindActiveByChatID(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if account != nil {
		shouldAskToDetach, currentUser, err := u.shouldAskToDetach(ctx, account)
		if err != nil {
			return err
		}
		if shouldAskToDetach {
			return u.Constructor.
				ConstructStartAttachedAndSend(ctx,
					[]int64{chatID},
					currentUser.GetFullNameWithSchool())
		}
	}

	if ticket != "" {
		shouldAskToAutoConnect, user, err := u.shouldAskToAutoConnect(ctx, ticket, now)
		if err != nil {
			return err
		}
		if shouldAskToAutoConnect {
			return u.ConnectTelegramToUserID(ctx, chatID, user.ID, user.SchoolName)
		}
	}

	if username != "" {
		user, err := u.repo.UsersRepository.ByUsername(ctx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if user.IsActivated && !user.IsDeleted {
			return u.ConnectTelegramToUserID(ctx, chatID, user.ID, user.SchoolName)
		}
	}

	return u.Constructor.
		ConstructStartManualAndSend(ctx,
			[]int64{chatID}, username)
}

func (u *Usecase) shouldAskToAutoConnect(ctx context.Context, token string, now time.Time) (bool, entity.User, error) {
	ticket, err := u.repo.TelegramAuthTicketsRepository.ByToken(ctx, token)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, entity.User{}, err
	}

	// Проверяем валидность тикета
	if ticket.IsSpent || ticket.IsBlocked || ticket.ExpiresAt.Before(now) {
		return false, entity.User{}, nil
	}

	user, err := u.repo.UsersRepository.ByID(ctx, ticket.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, entity.User{}, err
	}

	if user.IsActivated && !user.IsDeleted {
		return true, user, nil
	}

	return false, entity.User{}, nil
}

func (u *Usecase) shouldAskToDetach(ctx context.Context, account *entity.TelegramAccount) (bool, entity.User, error) {
	currentUser, err := u.repo.UsersRepository.ByID(ctx, account.UserID)
	if err != nil {
		fmt.Println(err)
		return false, entity.User{}, err
	}

	return currentUser.IsActivated && !currentUser.IsDeleted, currentUser, nil
}
