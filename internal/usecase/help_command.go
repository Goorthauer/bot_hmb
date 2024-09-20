package usecase

import (
	"context"
	"errors"
)

func (u *Usecase) HelpCommand(ctx context.Context, chatID int64) error {
	user, err := u.getUser(ctx, chatID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	switch {
	case user.IsMaster:
		return u.Constructor.ConstructMasterHelpAndSend(ctx,
			[]int64{chatID})
	case user.IsActivated:
		return u.Constructor.ConstructUserHelpAndSend(ctx,
			[]int64{chatID})
	default:
		return u.Constructor.ConstructHelpAndSend(ctx,
			[]int64{chatID})
	}
}
