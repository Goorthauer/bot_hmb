package usecase

import (
	"context"

	"github.com/go-telegram/bot/models"
)

func (u *Usecase) UnknownCommand(ctx context.Context, chatID int64, text string, update *models.Update) error {
	step, err := u.repo.StepRepository.ByChatID(ctx, chatID)
	if err != nil {
		return nil
	}
	if step.Step.IsValid() {
		if text == "" && update.Message != nil && update.Message.Contact != nil {
			text = update.Message.Contact.PhoneNumber
		}
		return u.RegistrationByStep(ctx, chatID, text)
	}
	err = u.Constructor.
		ConstructUnknownAndSendWithText(ctx,
			[]int64{chatID},
			text)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) UnknownErrCommand(ctx context.Context, chatID int64, e error) error {
	err := u.Constructor.
		ConstructUnknownErrAndSend(ctx,
			[]int64{chatID}, e)
	if err != nil {
		return err
	}

	return nil
}
