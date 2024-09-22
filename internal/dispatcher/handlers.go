package dispatcher

import (
	"bot_hmb/internal/entity"
	"context"
	"strings"

	"github.com/go-telegram/bot/models"
)

type Handlers struct {
	MasterHandlers           map[string]func(context.Context, int64, string) error
	MasterHandlersWithModels map[string]func(context.Context, int64, *models.Update) error
	UserHandlers             map[string]func(context.Context, int64, *models.Update) error
}

func (d *dispatcher) SetHandlers() {
	d.handlers = &Handlers{
		MasterHandlers: map[string]func(context.Context, int64, string) error{
			"/set_desc_and_price_schedule": d.setDescPriceSchedule,
			"/set_schedule":                d.setSchedule,
			"/set_subscriptions":           d.setSubscriptions,
			"/set_master":                  d.setMaster,
			"/create_ticket":               d.createTicket,
			"/hard_invite":                 d.hardInvite,
			"/subscription_list":           d.subscriptionList,
			"/subscription_quiz":           d.subscriptionQuiz,
		},
		MasterHandlersWithModels: map[string]func(context.Context, int64, *models.Update) error{
			"/present": d.subscriptionPresent,
		},
		UserHandlers: map[string]func(context.Context, int64, *models.Update) error{
			"/start":         d.startCommand,
			"/register":      d.registerCommand,
			"/subscriptions": d.subscriptionsCommand,
			"/detach":        d.detachCommand,
			"/info":          d.infoCommand,
			"/next_training": d.nextTrainingCommand,
			"/schedule":      d.scheduleCommand,
			"/help":          d.helpCommand,
		},
	}
}

func (d *dispatcher) checkHandler(ctx context.Context, chatID int64, text string, update *models.Update) error {
	if activeUser, err := d.userHandler(ctx, chatID, text, update); err != nil {
		return d.uc.UnknownErrCommand(ctx, chatID, err)
	} else if !activeUser {
		if activeMaster, err := d.masterHandler(ctx, chatID, text, update); err != nil {
			return d.uc.UnknownErrCommand(ctx, chatID, err)
		} else if !activeMaster {
			return d.uc.UnknownCommand(ctx, chatID, text, update)
		}
	}
	return nil
}

func (d *dispatcher) masterHandler(ctx context.Context, chatID int64, text string, update *models.Update) (bool, error) {
	for cmd, handler := range d.handlers.MasterHandlers {
		if strings.HasPrefix(text, cmd) {
			return true, handler(ctx, chatID, text)
		}
	}
	if update.CallbackQuery != nil {
		for cmd, handler := range d.handlers.MasterHandlers {
			if strings.HasPrefix(update.CallbackQuery.Data, cmd) {
				return true, handler(ctx, chatID, text)
			}
		}
		for cmd, handler := range d.handlers.MasterHandlersWithModels {
			if strings.HasPrefix(update.CallbackQuery.Data, cmd) {
				return true, handler(ctx, chatID, update)
			}
		}
	}
	return false, nil
}

func (d *dispatcher) userHandler(ctx context.Context, chatID int64, text string, update *models.Update) (bool, error) {
	if update.MyChatMember != nil && update.MyChatMember.NewChatMember.Type == models.ChatMemberTypeBanned {
		return true, d.uc.BannedCommand(ctx, chatID)
	}

	if update.CallbackQuery != nil {
		for cmd, handler := range d.handlers.UserHandlers {
			if strings.HasPrefix(update.CallbackQuery.Data, cmd) {
				return true, handler(ctx, chatID, update)
			}
		}
	}

	for cmd, handler := range d.handlers.UserHandlers {
		if strings.HasPrefix(text, cmd) {
			return true, handler(ctx, chatID, update)
		}
	}

	return false, nil
}

// Команды мастера
func (d *dispatcher) subscriptionQuiz(ctx context.Context, chatID int64, _ string) error {
	return d.uc.SubscriptionQuiz(ctx, chatID)
}

func (d *dispatcher) subscriptionList(ctx context.Context, chatID int64, _ string) error {
	return d.uc.GetSubscriptionList(ctx, chatID)
}

func (d *dispatcher) setDescPriceSchedule(ctx context.Context, chatID int64, text string) error {
	args := parseArgs(text, "/set_desc_and_price_schedule", ";")
	price := ""
	desc := ""
	if len(args) > 1 {
		price = strings.TrimSpace(args[0])
		desc = strings.TrimSpace(args[1])
	}
	return d.uc.SetDescPriceSchedule(ctx, chatID, price, desc)
}

func (d *dispatcher) setSchedule(ctx context.Context, chatID int64, text string) error {
	args := parseArgs(text, "/set_schedule", ";")
	scheduleList := make([]entity.TrainingDay, 0)
	for _, v := range args {
		dataDays := strings.Split(v, ":")
		if len(dataDays) > 2 {
			scheduleList = append(scheduleList, entity.TrainingDay{
				Day: entity.TrainingDayDay(strings.TrimSpace(dataDays[0])),
				Time: entity.TrainingDayTime{
					Open:   strings.TrimSpace(dataDays[1]),
					Closed: strings.TrimSpace(dataDays[2]),
				},
				Description: strings.TrimSpace(dataDays[3]),
			})
		}
	}
	return d.uc.SetSchedule(ctx, chatID, scheduleList)
}

func (d *dispatcher) setMaster(ctx context.Context, chatID int64, text string) error {
	args := strings.Split(text, " ")
	phone := ""
	if len(args) > 1 {
		phone = args[1]
	}
	return d.uc.SetMasterUser(ctx, chatID, phone)
}

func (d *dispatcher) setSubscriptions(ctx context.Context, chatID int64, text string) error {
	args := strings.Split(text, " ")
	var phone, days, price string
	if len(args) > 2 {
		phone = args[1]
		days = args[2]
	}
	if len(args) > 3 {
		price = args[3]
	}
	return d.uc.SetSubscriptions(ctx, chatID, days, phone, price)
}

func (d *dispatcher) subscriptionPresent(ctx context.Context, chatID int64, update *models.Update) error {
	data := strings.Split(update.CallbackQuery.Data, " ")
	userID := ""
	if len(data) > 1 {
		userID = data[1]
	}
	return d.uc.SubscriptionPresent(ctx, chatID, update.CallbackQuery.Message.Message.ID, userID)
}

func (d *dispatcher) createTicket(ctx context.Context, chatID int64, text string) error {
	args := strings.Split(text, " ")
	username := ""
	if len(args) > 1 {
		username = args[1]
	}
	return d.uc.CreateTelegramAuthTicket(ctx, chatID, username)
}

func (d *dispatcher) hardInvite(ctx context.Context, chatID int64, text string) error {
	args := strings.Split(text, " ")
	var username, firstName, lastName, phone, schoolID string
	if len(args) > 4 {
		username = args[1]
		phone = args[2]
		firstName = args[3]
		lastName = args[4]
	}
	if len(args) > 5 {
		schoolID = args[5]
	}
	return d.uc.InviteUser(ctx, chatID, schoolID, username, phone, firstName, lastName)
}

// Команды пользователя

func (d *dispatcher) startCommand(ctx context.Context, chatID int64, update *models.Update) error {
	ticket := ""
	if update.CallbackQuery != nil {
		ticket = strings.TrimSpace(update.CallbackQuery.Data[len("/start"):])
	} else {
		tokens := strings.Split(update.Message.Text, " ")
		if len(tokens) > 1 {
			ticket = tokens[1]
		}
	}
	username := update.Message.Chat.Username
	if update.CallbackQuery != nil {
		username = update.CallbackQuery.Message.Message.Chat.Username
	}
	return d.uc.StartCommand(ctx, chatID, ticket, username)
}

func (d *dispatcher) registerCommand(ctx context.Context, chatID int64, update *models.Update) error {
	return d.uc.Registration(ctx, chatID, update.Message.Chat.Username)
}

func (d *dispatcher) registerLastStep(ctx context.Context, chatID int64, update *models.Update) error {
	data := strings.Split(update.CallbackQuery.Data, " ")
	schoolID := ""
	if len(data) > 1 {
		schoolID = data[1]
	}
	return d.uc.RegistrationLastStep(ctx,
		chatID, schoolID,
	)
}

func (d *dispatcher) subscriptionsCommand(ctx context.Context, chatID int64, _ *models.Update) error {
	return d.uc.GetSubscriptions(ctx, chatID)
}

func (d *dispatcher) detachCommand(ctx context.Context, chatID int64, _ *models.Update) error {
	return d.uc.DetachCommand(ctx, chatID)
}

func (d *dispatcher) infoCommand(ctx context.Context, chatID int64, _ *models.Update) error {
	return d.uc.InfoCommand(ctx, chatID)
}

func (d *dispatcher) nextTrainingCommand(ctx context.Context, chatID int64, _ *models.Update) error {
	return d.uc.NextTrainingDay(ctx, chatID)
}

func (d *dispatcher) scheduleCommand(ctx context.Context, chatID int64, _ *models.Update) error {
	return d.uc.GetSchedule(ctx, chatID)
}

func (d *dispatcher) helpCommand(ctx context.Context, chatID int64, _ *models.Update) error {
	return d.uc.HelpCommand(ctx, chatID)
}

func parseArgs(text, prefix, sep string) []string {
	text = strings.TrimPrefix(text, prefix)
	return strings.Split(text, sep)
}
