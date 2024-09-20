package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func (u *Usecase) BannedCommand(ctx context.Context, chatID int64) error {
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
	}

	return nil
}
