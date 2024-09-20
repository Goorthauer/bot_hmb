package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/gofrs/uuid"
)

func (u *Usecase) CreateTelegramAuthTicket(
	ctx context.Context,
	chatID int64,
	username string) error {
	masterUser, err := u.getUser(ctx, chatID)
	if err != nil {
		return err
	}
	if !masterUser.IsMaster {
		return ErrIsNotMaster
	}

	user, err := u.repo.UsersRepository.ByUsername(ctx, username)
	if err != nil {
		return err
	}
	url, err := u.createTicket(ctx, user.ID)
	u.Wrapper.SendMessage(ctx, bot.EscapeMarkdown(url), strconv.Itoa(int(chatID)))
	return nil
}

func randomString(length int) (string, error) {
	buff := make([]byte, length)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(buff)
	return str[:length], nil
}

func (u *Usecase) createTicket(ctx context.Context, userID uuid.UUID) (string, error) {
	now := time.Now().UTC()
	token, err := randomString(64)
	if err != nil {
		return "", err
	}

	ticket := entity.TelegramAuthTicket{
		Token:     token,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(48 * time.Hour),
	}

	oldTickets, err := u.repo.TelegramAuthTicketsRepository.FindActiveTickets(ctx, userID)
	if err != nil {
		return "", err
	}

	err = u.repo.TelegramAuthTicketsRepository.Upsert(ctx, &ticket)
	if err != nil {
		return "", err
	}

	for _, oldTicket := range oldTickets {
		err = u.repo.TelegramAuthTicketsRepository.
			DeactivateTicket(ctx, oldTicket.Token)
		if err != nil {
			return "", err
		}
	}

	url := fmt.Sprintf("%s?start=%s", u.config.TelegramBotURL, ticket.Token)
	return url, nil
}
