package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/gofrs/uuid"
)

func (u *Usecase) SubscriptionPresent(ctx context.Context,
	chatID int64, messageID int, userID string) error {
	masterUser, err := u.getUser(ctx, chatID)
	if errors.Is(err, ErrNotFound) {
		return u.Constructor.
			ConstructUnknownAndSend(ctx,
				[]int64{chatID})
	}
	if err != nil {
		return err
	}
	if !masterUser.IsMaster {
		return ErrIsNotMaster
	}
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return err
	}
	user, err := u.repo.UsersRepository.ByID(ctx, userUUID)
	if err != nil {
		return err
	}
	presentID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	err = u.repo.UserPresentsRepository.Create(ctx, &entity.Presents{
		ID:        presentID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	text := fmt.Sprintf(`%s отмечен как присутствующий.`, user.PersonalData.GetFullName())
	u.Wrapper.EditMessageText(ctx,
		bot.EscapeMarkdown(text),
		strconv.Itoa(int(chatID)), messageID)
	return nil
}
