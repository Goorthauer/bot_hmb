package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"errors"

	"github.com/go-telegram/bot"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func (u *Usecase) ConnectTelegramToUserID(ctx context.Context,
	chatID int64, userID uuid.UUID,
	schoolName string) error {
	var accountFound = false
	oldAccount, err := u.repo.TelegramAccountsRepository.FindActiveByChatID(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if err == nil {
		accountFound = true
	}

	if accountFound && oldAccount.UserID != userID {
		err = u.detachTelegramAndSendMessage(ctx, oldAccount.ChatID, oldAccount.UserID)
		if err != nil {
			return err
		}
	}

	if !accountFound || oldAccount.UserID != userID {
		err = u.repo.TelegramAccountsRepository.
			Upsert(ctx, &entity.TelegramAccount{
				UserID:   userID,
				ChatID:   chatID,
				IsActive: true,
			})
		if err != nil {
			return err
		}

		err = u.Constructor.ConstructAttachedAndSend(ctx,
			[]int64{chatID}, bot.EscapeMarkdown(schoolName))
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) detachTelegramAndSendMessage(ctx context.Context,
	chatID int64, userID uuid.UUID) error {
	err := u.repo.TelegramAccountsRepository.DetachChatID(ctx, chatID, userID)
	if err != nil {
		return err
	}

	err = u.Constructor.ConstructDetachedAndSend(ctx,
		[]int64{chatID})
	if err != nil {
		return err
	}

	return nil
}
