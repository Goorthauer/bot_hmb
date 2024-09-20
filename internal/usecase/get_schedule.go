package usecase

import (
	"context"
	"errors"
)

func (u *Usecase) GetSchedule(ctx context.Context, chatID int64) error {
	user, err := u.getUser(ctx, chatID)
	if errors.Is(err, ErrNotFound) {
		return u.Constructor.
			ConstructUnknownAndSend(ctx,
				[]int64{chatID})
	}
	if err != nil {
		return err
	}
	training, err := u.repo.SchoolsTrainingRepository.BySchool(ctx, user.SchoolID)
	if err != nil {
		return err
	}
	return u.Constructor.ConstructScheduleAndSend(ctx,
		[]int64{chatID},
		user.PersonalData.GetFullName(),
		user.SchoolName,
		training.Price,
		training.Description,
		training.Schedule)

}
