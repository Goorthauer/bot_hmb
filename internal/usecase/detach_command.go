package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func (u *Usecase) DetachCommand(ctx context.Context, chatID int64) error {
	var accountFound = false
	account, err := u.repo.TelegramAccountsRepository.FindActiveByChatID(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if err == nil {
		accountFound = true
	}

	if accountFound {
		err = u.repo.TelegramAccountsRepository.DetachChatID(ctx, account.ChatID, account.UserID)
		if err != nil {
			return err
		}

		err = u.Constructor.
			ConstructDetachedAndSend(ctx, []int64{account.ChatID})
		if err != nil {
			return err
		}

		return nil
	}

	err = u.Constructor.
		ConstructDetachDetachedAndSend(ctx,
			[]int64{chatID})
	if err != nil {
		return err
	}

	return nil
}
