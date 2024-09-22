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
	"gorm.io/gorm"
)

func (u *Usecase) SetSubscriptions(ctx context.Context,
	chatID int64,
	days, phone, price string) error {
	const countDayOnWeek = 7.0
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
		return nil
	}
	daysNumber, err := strconv.Atoi(days)
	if err != nil {
		return err
	}
	subID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	user, err := u.repo.UsersRepository.ByPhone(ctx, phone)
	if err != nil {
		return err
	}
	now := time.Now()
	oldSub, err := u.repo.SubscriptionsRepository.ByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	baseDeadline := now
	if time.Since(oldSub.DeadlineAt) < 0 {
		baseDeadline = oldSub.DeadlineAt
	}
	trainingSchedule, err := u.repo.SchoolsTrainingRepository.BySchool(ctx, user.SchoolID)
	if err != nil {
		return err
	}
	if price == "" {
		price = strconv.Itoa(trainingSchedule.Price)
	}
	if len(trainingSchedule.Schedule) <= 0 || len(trainingSchedule.Schedule) > countDayOnWeek {
		u.Wrapper.SendMessage(ctx, bot.EscapeMarkdown(`у данной школы ошибка в расписании`), strconv.Itoa(int(chatID)))
		return nil
	}
	trainingDayPerWeek := countDayOnWeek / float64(len(trainingSchedule.Schedule))
	var additionalDays = 0
	switch {
	case daysNumber > 10:
		additionalDays = 2
	case daysNumber > 6:
		additionalDays = 1
	}
	addDays := int(trainingDayPerWeek*float64(daysNumber)) + additionalDays
	deadline := baseDeadline.AddDate(0, 0, addDays)
	err = u.repo.SubscriptionsRepository.Create(ctx, &entity.Subscription{
		ID:         subID,
		UserID:     user.ID,
		SchoolID:   user.SchoolID,
		CreatedAt:  now,
		DeadlineAt: deadline,
		Price:      price,
		Days:       daysNumber,
	})
	if err != nil {
		return err
	}
	text := fmt.Sprintf(

		`Абонемент пользователю '%s' успешно выдан.
	Он действует до: %v
`,
		user.PersonalData.GetFullName(), deadline.Format(time.DateOnly))
	u.Wrapper.SendMessage(ctx, bot.EscapeMarkdown(text), strconv.Itoa(int(chatID)))
	return nil

}
