package dispatcher

import (
	"bot_hmb/internal/usecase"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, _ *bot.Bot, update *models.Update)
	SetHandlers()
}

type dispatcher struct {
	uc       *usecase.Usecase
	handlers *Handlers
}

func NewDispatcher(
	uc *usecase.Usecase) Dispatcher {
	return &dispatcher{
		uc: uc,
	}
}

func (d *dispatcher) Dispatch(ctx context.Context, _ *bot.Bot, update *models.Update) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic recovered in telegram bot")
		}
	}()

	if err := d.dispatch(ctx, update); err != nil {
		fmt.Println("error in telegram bot: ", err.Error())
	}
}

func (d *dispatcher) dispatch(ctx context.Context, update *models.Update) error {
	chatID := d.getChatID(update)
	if chatID == -1 {
		return errors.New("could not find chat id")
	}

	text := ""
	if update.Message != nil {
		text = update.Message.Text
		fmt.Printf(">###\ntext command:%v\nfrom chat:%v\n<###\n", text, chatID)
	}
	if update.CallbackQuery != nil {
		fmt.Printf(">###\ncallback data:%v\nfrom chat:%v\n<###\n", update.CallbackQuery.Data, chatID)
	}
	data := strings.Split(text, "\n")
	for _, v := range data {
		err := d.checkHandler(ctx, chatID, v, update)
		if err != nil {
			return fmt.Errorf("check handler err %w", err)
		}
	}
	return nil
}

func (d *dispatcher) getChatID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.MyChatMember != nil {
		return update.MyChatMember.Chat.ID
	}

	if update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}

	return -1
}
