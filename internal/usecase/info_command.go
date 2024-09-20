package usecase

import (
	"context"
	"errors"
)

func (u *Usecase) InfoCommand(ctx context.Context, chatID int64) error {
	user, err := u.getUser(ctx, chatID)
	if errors.Is(err, ErrNotFound) {
		return u.Constructor.
			ConstructInfoDetachedAndSend(ctx,
				[]int64{chatID})
	}
	if err != nil {
		return err
	}

	return u.Constructor.
		ConstructInfoAttachedAndSend(ctx,
			[]int64{chatID},
			user.GetFullNameWithSchool())
}
