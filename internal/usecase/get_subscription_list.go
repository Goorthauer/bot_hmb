package usecase

import (
	"bot_hmb/internal/entity"
	"context"
	"errors"

	"github.com/gofrs/uuid"
)

func (u *Usecase) GetSubscriptionList(ctx context.Context, chatID int64) error {
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
	subs, err := u.repo.UserPresentsRepository.BySchoolIDWithLost(ctx, masterUser.SchoolID)
	if err != nil {
		return err
	}
	userIDs := make([]uuid.UUID, 0, len(subs))
	for _, v := range subs {
		userIDs = append(userIDs, v.UserID)
	}
	users, err := u.repo.UsersRepository.ByIDs(ctx, userIDs)
	if err != nil {
		return err
	}
	userList := make(map[uuid.UUID]entity.User)
	for _, v := range users {
		if v == nil {
			continue
		}
		userList[v.ID] = *v
	}
	return u.Constructor.ConstructSubscriptionListAndSendV2(ctx,
		[]int64{chatID},
		masterUser.GetFullNameWithSchool(),
		userList,
		subs)

}
