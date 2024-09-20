package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
)

func (u *Usecase) SetMasterUser(
	ctx context.Context,
	chatID int64,
	phone string) error {
	masterUser, err := u.getMaster(ctx, chatID)
	if err != nil {
		return err
	}
	if u.config.MasterUserNickname == masterUser.Username && phone == "" {
		phone = masterUser.Phone
	}
	if _, err := checkPhone(phone); err != nil {
		return err
	}
	user, err := u.repo.UsersRepository.ByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if !user.IsMaster {
		err = u.repo.UsersRepository.SetMasterRights(ctx, user.ID, true)
		if err != nil {
			return err
		}
	}
	text := fmt.Sprintf(
		`Права пользователю '%s' успешно выданы.Теперь он мастер-юзер.`,
		user.PersonalData.GetFullName())
	u.Wrapper.SendMessage(ctx, bot.EscapeMarkdown(text), strconv.Itoa(int(chatID)))
	return nil
}
