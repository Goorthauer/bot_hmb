package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func (u *Usecase) GetSubscriptions(ctx context.Context, chatID int64) error {
	user, err := u.getUser(ctx, chatID)
	if errors.Is(err, ErrNotFound) {
		return u.Constructor.
			ConstructUnknownAndSend(ctx,
				[]int64{chatID})
	}
	if err != nil {
		return err
	}
	sub, err := u.repo.SubscriptionsRepository.ByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.Constructor.ConstructSubscriptionsFailAndSend(ctx, []int64{chatID}, user.GetFullNameWithSchool())
		}
		return err
	}
	return u.Constructor.ConstructSubscriptionsAndSend(ctx,
		[]int64{chatID},
		user.PersonalData.GetFullName(),
		sub.DeadlineAt)

}
