package usecase

import (
	"context"
	"errors"
)

func (u *Usecase) SubscriptionQuiz(ctx context.Context, chatID int64) error {
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
		return ErrIsNotMaster
	}
	users, err := u.repo.UsersRepository.BySchool(ctx, masterUser.SchoolID, true)
	if err != nil {
		return err
	}
	return u.Constructor.ConstructSubscriptionQuizAndSend(ctx,
		[]int64{chatID},
		masterUser.GetFullNameWithSchool(),
		users)
}
