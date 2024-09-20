package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
)

func (u *Usecase) SetDescPriceSchedule(ctx context.Context,
	chatID int64,
	rawPrice, desc string) error {
	masterUser, err := u.getUser(ctx, chatID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return u.Constructor.
				ConstructUnknownAndSend(ctx,
					[]int64{chatID})
		}
		return err
	}
	if !masterUser.IsMaster {
		return nil
	}
	if rawPrice == "" && desc == "" {
		return nil
	}
	schedule, err := u.repo.SchoolsTrainingRepository.BySchool(ctx, masterUser.SchoolID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			text := fmt.Sprintf(
				`Сначала необходимо создать расписание занятий.`)
			u.Wrapper.SendMessage(ctx, bot.EscapeMarkdown(text), strconv.Itoa(int(chatID)))

		}
		return err
	}
	if rawPrice != "" {
		price, err := strconv.Atoi(rawPrice)
		if err != nil {
			return err
		}
		schedule.Price = price
	}
	if desc != "" {
		schedule.Description = desc
	}
	err = u.repo.SchoolsTrainingRepository.Create(ctx, &schedule)
	if err != nil {
		return err
	}
	return u.Constructor.ConstructScheduleAndSend(ctx,
		[]int64{chatID},
		masterUser.PersonalData.GetFullName(),
		masterUser.SchoolName,
		schedule.Price,
		schedule.Description,
		schedule.Schedule)
}
