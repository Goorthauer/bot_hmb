package usecase

import (
	"context"
	"errors"
)

func (u *Usecase) NextTrainingDay(ctx context.Context, chatID int64) error {
	user, err := u.getUser(ctx, chatID)
	if errors.Is(err, ErrNotFound) {
		return u.Constructor.
			ConstructInfoDetachedAndSend(ctx,
				[]int64{chatID})
	}
	if err != nil {
		return err
	}
	training, err := u.repo.SchoolsTrainingRepository.BySchool(ctx, user.SchoolID)
	if err != nil {
		return err
	}
	dateNextTraining, t, err := training.FindNextTrainingDay()
	if err != nil {
		return err
	}
	return u.Constructor.
		ConstructNextTrainingDayAndSend(ctx,
			[]int64{chatID},
			user.GetFullNameWithSchool(),
			user.SchoolAddress,
			t,
			dateNextTraining)
}
