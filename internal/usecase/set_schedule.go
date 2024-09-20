package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"errors"
	"time"
)

func (u *Usecase) SetSchedule(ctx context.Context,
	chatID int64,
	schedule []entity.TrainingDay) error {
	const price = 2500
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
	outSchedule := make([]entity.TrainingDay, 0, len(schedule))
	for _, v := range schedule {
		outSchedule = append(outSchedule, v)
	}
	err = u.repo.SchoolsTrainingRepository.Create(ctx, &entity.SchoolTraining{
		SchoolID:  masterUser.SchoolID,
		Schedule:  outSchedule,
		Price:     price,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return u.Constructor.ConstructScheduleAndSend(ctx,
		[]int64{chatID},
		masterUser.PersonalData.GetFullName(),
		masterUser.SchoolName,
		price,
		"",
		outSchedule)
}
