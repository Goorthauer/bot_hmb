package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
)

func validateUserData(firstName, lastName, phone string) error {
	if firstName == "" || lastName == "" || phone == "" {
		return errors.New("firstName, lastName или phone не могут быть пустыми")
	}
	if _, err := checkPhone(phone); err != nil {
		return err
	}
	return nil
}

func (u *Usecase) InviteUser(
	ctx context.Context,
	chatID int64,
	schoolID,
	username,
	phone,
	firstName,
	lastName string) error {

	if err := validateUserData(firstName, lastName, phone); err != nil {
		return err
	}

	masterUser, err := u.getUser(ctx, chatID)
	if err != nil {
		return err
	}
	if !masterUser.IsMaster {
		return ErrIsNotMaster
	}
	if schoolID == "" {
		schoolID = masterUser.SchoolID.String()
	}

	user, err := u.registerUser(ctx, schoolID, firstName, lastName, phone, username)
	if err != nil {
		return err
	}

	url, err := u.createTicket(ctx, user.ID)
	if err != nil {
		return err
	}

	message := fmt.Sprintf(
		"пользователь %s %s успешно создан, теперь ему можно прислать приглашение: %s",
		bot.EscapeMarkdown(firstName),
		bot.EscapeMarkdown(lastName),
		url,
	)
	u.Wrapper.SendMessage(ctx, message, strconv.Itoa(int(chatID)))

	return nil
}
