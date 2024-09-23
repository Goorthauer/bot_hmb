package telegram

import (
	"bot_hmb/internal/entity"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/gofrs/uuid"
)

type Constructor interface {
	ConstructUnknownErrAndSend(ctx context.Context, chatIDs []int64, e error) error
	ConstructUnknownAndSend(ctx context.Context, chatIDs []int64) error
	ConstructUnknownAndSendWithText(ctx context.Context, chatIDs []int64, text string) error
	ConstructInfoAttachedAndSend(ctx context.Context, chatIDs []int64, fullname string) error
	ConstructNextTrainingDayAndSend(ctx context.Context, chatIDs []int64, fioFull, schoolAddress, trainingTime string, date time.Time) error
	ConstructHelpAndSend(ctx context.Context, chatIDs []int64) error
	ConstructMasterHelpAndSend(ctx context.Context, chatIDs []int64) error
	ConstructStartAttachedAndSend(ctx context.Context, chatIDs []int64, fioFull string) error
	ConstructStartAutoAndSend(ctx context.Context, chatIDs []int64, username string) error
	ConstructStartManualAndSend(ctx context.Context, chatIDs []int64, username string) error
	ConstructInfoDetachedAndSend(ctx context.Context, chatIDs []int64) error
	ConstructDetachedAndSend(ctx context.Context, chatIDs []int64) error
	ConstructDetachDetachedAndSend(ctx context.Context, chatIDs []int64) error
	ConstructAttachedAndSend(ctx context.Context, chatIDs []int64, schoolName string) error
	ConstructSubscriptionsAndSend(ctx context.Context, chatIDs []int64, fioFull string, deadline time.Time) error
	ConstructSubscriptionListAndSend(ctx context.Context, chatIDs []int64, fioFull string, userList map[uuid.UUID]entity.User, subs []entity.Subscription) error
	ConstructSubscriptionListAndSendV2(ctx context.Context, chatIDs []int64, fioFull string, userList map[uuid.UUID]entity.User, subs []*entity.PresentsSubscription) error
	ConstructSubscriptionQuizAndSend(ctx context.Context, chatIDs []int64, fioFull string, users []*entity.User) error
	ConstructUserHelpAndSend(ctx context.Context, chatIDs []int64) error
	ConstructRegisterPhone(ctx context.Context, chatID int64) error
	ConstructScheduleAndSend(ctx context.Context, chatIDs []int64, fioFull string, schoolName string, price int, desc string, schedule entity.Schedule) error
	ConstructRegisterWithSchool(ctx context.Context, chatID int64, schools map[uuid.UUID]string) error
	ConstructRegisterWithStep(ctx context.Context, chatID int64, step entity.StepType, data map[uuid.UUID]string) error
	ConstructSubscriptionsFailAndSend(ctx context.Context, chatIDs []int64, fioFull string) error
}

type constructor struct {
	debug  bool
	sender Wrapper
}

func NewConstructor(debug bool, wrapper Wrapper) Constructor {
	return &constructor{debug: debug, sender: wrapper}
}

func (c *constructor) ConstructInfoDetachedAndSend(ctx context.Context, chatIDs []int64) error {
	text := "ℹ️ На текущий момент бот не привязан к учётной записи HMB Schools\\.\n\n" +
		"Для подключения нужно войти в свою учётную запись 👇"

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Войти",
					CallbackData: "/start",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructUnknownAndSend(ctx context.Context, chatIDs []int64) error {
	text := "🤔 Я не знаю, что это означает\\. Но могу показать справку :\\)\n\n" +
		"или введите команду самостоятельно \\- \\/help"

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Справка",
					CallbackData: "/help",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructUnknownAndSendWithText(ctx context.Context, chatIDs []int64, text string) error {
	text = fmt.Sprintf("🤔 Я не знаю, что означает команда `%s`\\. Но могу показать справку :\\)\n\n"+
		"или введите команду самостоятельно \\- \\/help", text)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Справка",
					CallbackData: "/help",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructUnknownErrAndSend(ctx context.Context, chatIDs []int64, e error) error {
	text := "⛔⛔⛔ Произошла ошибка\\. Проверьте правильность ввода\n\\."

	if c.debug {
		text += bot.EscapeMarkdown(e.Error())
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Справка",
					CallbackData: "/help",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructInfoAttachedAndSend(ctx context.Context, chatIDs []int64, fullname string) error {
	text := fmt.Sprintf(
		"ℹ️ Бот подключен к учётной записи *%s*"+
			"\n\n"+
			"Чтобы отключить аккаунт — используйте команду /detach или нажмите на кнопку ниже 👇",
		fullname,
	)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         fmt.Sprintf("Выйти из %s", fullname),
					CallbackData: "/detach",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, bot.EscapeMarkdown(text), kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructSubscriptionListAndSend(ctx context.Context, chatIDs []int64, fioFull string, userList map[uuid.UUID]entity.User, subs []entity.Subscription) error {
	text := fmt.Sprintf(
		`%s, вот список людей вашей школы:
`, fioFull)
	if len(subs) == 0 || len(userList) == 0 {
		text = fmt.Sprintf(
			`%s, К сожалению нет зарегистрированных абонементов.
что бы зарегистрировать абонемент - осуществите команду:
/set_subscriptions {номер телефона} {количество дней} {*цена}
`, fioFull)
	}
	for _, v := range subs {
		var user entity.User
		var ok bool
		if user, ok = userList[v.UserID]; !ok {
			continue
		}
		danger := ""
		dur := time.Since(v.DeadlineAt)
		dayExp := math.Round(-dur.Hours() / 24)
		switch {
		case dayExp == 0:
			danger = "(⚠️закончится сегодня) "
		case dayExp < 0:
			danger = "‼️закончился - "
		case dayExp < 3:
			danger = fmt.Sprintf(`(закончится через %v дней) `, dayExp)
		}

		text += fmt.Sprintf(
			`
💪 %s(%s) - абонемент до %s%s;
`,
			user.PersonalData.GetFullName(), user.Phone, danger, v.DeadlineAt.Format(time.DateOnly))
	}
	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructSubscriptionListAndSendV2(ctx context.Context, chatIDs []int64, fioFull string, userList map[uuid.UUID]entity.User, subs []*entity.PresentsSubscription) error {
	text := fmt.Sprintf(
		`%s, вот список людей вашей школы:
`, fioFull)
	if len(subs) == 0 || len(userList) == 0 {
		text = fmt.Sprintf(
			`%s, К сожалению нет зарегистрированных абонементов.
что бы зарегистрировать абонемент - осуществите команду:
/set_subscriptions {номер телефона} {количество дней} {*цена}
`, fioFull)
	}
	for _, v := range subs {
		if v == nil {
			continue
		}
		var user entity.User
		var ok bool
		if user, ok = userList[v.UserID]; !ok {
			continue
		}
		danger := ""
		if v.LostDays >= v.CountTraining {
			dur := time.Since(v.DeadlineAt)
			dayExp := math.Round(-dur.Hours() / 24)
			switch {
			case dayExp == 0:
				danger = "(⚠️закончится сегодня) "
			case dayExp < 0:
				danger = "‼️закончился - "
			case dayExp < 3:
				danger = fmt.Sprintf(`(закончится через %v дней) `, dayExp)
			}
		} else {
			danger = "‼️закончился, т.к. не осталось занятий "
		}

		text += fmt.Sprintf(
			`
💪 %s(%s) - абонемент до %s%s(осталось %v занятий);
`,
			user.PersonalData.GetFullName(), user.Phone, danger, v.DeadlineAt.Format(time.DateOnly),
			v.LostDays-v.CountTraining)
	}
	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructSubscriptionQuizAndSend(ctx context.Context, chatIDs []int64, fioFull string, users []*entity.User) error {
	text := fmt.Sprintf("%s, Отметь людей присутствующих сегодня на занятии.", fioFull)
	if len(users) != 0 {
		err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
		if err != nil {
			return err
		}
	}
	for _, user := range users {
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         user.PersonalData.GetFullName(),
						CallbackData: fmt.Sprintf("/present %s", user.ID),
					},
				},
			},
		}
		err := c.sendMessageWithButtonsToChatIDs(ctx, "был ли человек на тренировке\\?", &kb, chatIDs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *constructor) ConstructSubscriptionsAndSend(ctx context.Context, chatIDs []int64, fioFull string, deadline time.Time) error {
	var text string
	if time.Since(deadline) < 0 {
		text = fmt.Sprintf(`🏷%s, ваш абонемент действует до %s`,
			fioFull,
			deadline.Format(time.DateOnly))
	} else {
		text = fmt.Sprintf(`🏷%s, ваш закончился %s`,
			fioFull,
			deadline.Format(time.DateOnly))
	}
	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructSubscriptionsFailAndSend(ctx context.Context, chatIDs []int64, fioFull string) error {
	text := fmt.Sprintf(`🏷%s, ваш абонемент еще не зарегистрирован в системе😭😭😭`,
		fioFull)
	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
	if err != nil {
		return err
	}

	return nil
}
func (c *constructor) ConstructNextTrainingDayAndSend(ctx context.Context, chatIDs []int64, fioFull, schoolAddress, trainingTime string, date time.Time) error {
	text := fmt.Sprintf(`🏷%s, дата следующей тренировки(по расписанию):
🕗%v, тренировка проходит %s по адресу %s`,
		fioFull, date.Format(time.DateOnly), trainingTime, schoolAddress)
	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructHelpAndSend(ctx context.Context, chatIDs []int64) error {
	err := c.sendMessageToChatIDs(ctx, getHelpBasicText(), chatIDs)
	if err != nil {
		return err
	}

	return nil
}
func (c *constructor) ConstructUserHelpAndSend(ctx context.Context, chatIDs []int64) error {
	text := getHelpBasicText() + bot.EscapeMarkdown(getHelpUserText())
	err := c.sendMessageToChatIDs(ctx, text, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructRegisterPhone(ctx context.Context, chatID int64) error {
	phoneButton := models.KeyboardButton{
		Text:           "Отправить номер телефона",
		RequestContact: true,
	}

	replyKeyboard := models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{phoneButton},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	text := "Пожалуйста, отправьте ваш номер телефона:\nКнопкой ниже или самостоятельно в формате 7ХХХХХХХХХХ"
	err := c.sendMessageWithButtonsToChatIDs(ctx, bot.EscapeMarkdown(text), replyKeyboard, []int64{chatID})
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructRegisterWithSchool(ctx context.Context, chatID int64, schools map[uuid.UUID]string) error {
	text := "Последний шаг! Выберите тренировочный зал, который ближе всего к вам находится."
	buttons := make([][]models.InlineKeyboardButton, 0)
	for id, v := range schools {
		callbackData := fmt.Sprintf("/last-step %s", id)
		buttons = append(buttons, []models.InlineKeyboardButton{{
			Text:         v,
			CallbackData: callbackData,
		}})
	}
	kb := &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
	err := c.sendMessageWithButtonsToChatIDs(ctx, bot.EscapeMarkdown(text), kb, []int64{chatID})
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructRegisterWithStep(ctx context.Context, chatID int64, step entity.StepType, data map[uuid.UUID]string) error {
	var text string
	switch step {
	case entity.StepFirstName:
		text = "Введите свое имя"
	case entity.StepLastname:
		text = "Введите свою фамилию"
	case entity.StepPhone:
		c.ConstructRegisterPhone(ctx, chatID)
		//text = "Введите номер своего телефона в формате 7ХХХХХХХХХХ"
	case entity.StepSchool:
		c.ConstructRegisterWithSchool(ctx, chatID, data)
	default:
		text = "Произошла неизвестная ошибка!"

	}

	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), []int64{chatID})
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructMasterHelpAndSend(ctx context.Context, chatIDs []int64) error {
	text := getHelpBasicText()

	text += bot.EscapeMarkdown(getHelpUserText() + getHelpMasterText())
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Узнать дату следующей тренировки",
					CallbackData: "/next_training",
				},
			},
			{
				{
					Text:         "Список абонементов",
					CallbackData: "/subscription_list",
				},
			},
			{
				{
					Text:         "Отметить людей на тренировке",
					CallbackData: "/subscription_quiz",
				},
			},

			{
				{
					Text:         "Актуальное расписание",
					CallbackData: "/schedule",
				},
			},
		},
	}
	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructStartAttachedAndSend(ctx context.Context,
	chatIDs []int64,
	fioFull string) error {
	text := fmt.Sprintf(
		"⚠️ Бот уже подключен к учётной записи *%s*, "+
			"если нужно подключить другую — сначала нужно отвязать аккаунт",
		bot.EscapeMarkdownUnescaped(fioFull))

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         fmt.Sprintf("Выйти из %s", fioFull),
					CallbackData: "/detach",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructStartAutoAndSend(ctx context.Context, chatIDs []int64, username string) error {
	text := fmt.Sprintf(
		"👋 Привет %s\\!"+
			"\n\n"+
			"Это бот HMB Schools\\. "+
			"Через него мы рассылаем важные уведомления\\."+
			"\n\n"+
			"🥑 Для подключения аккаунта *%s* нажмите на кнопку ниже 👇"+
			"\n\n"+
			"_Если в процессе возникнут ошибки или нужен другой аккаунт – используйте подключение вручную, "+
			"для этого достаточно войти в свою учётную запись_",
		username,
		username)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Регистрация",
					CallbackData: "/register",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructStartManualAndSend(ctx context.Context, chatIDs []int64, username string) error {
	text := fmt.Sprintf(
		"👋 Привет, %s\\!"+
			"\n\n"+
			"Это бот HMB Schools\\. "+
			"Через него мы рассылаем важные уведомления\\."+
			"\n\n"+
			"🐒Если вы хотите самостоятельно зарегистрироваться, тогда используйте команду"+
			"\n\n"+
			"/register"+
			"\n\n"+
			"🐣Если вы уже зарегстрированы тренером, то просто нажмите кнопку ниже 👇",
		username)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Войти",
					CallbackData: "/start",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructDetachedAndSend(ctx context.Context, chatIDs []int64) error {
	text := "🖇 Отвязали бота от вашей учётной записи, теперь уведомления приходить *не будут*" +
		"\n\n" +
		"Если это произошло по ошибке — войдите в приложение ещё раз 👇"

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Войти",
					CallbackData: "/start",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructDetachDetachedAndSend(ctx context.Context, chatIDs []int64) error {
	text := "👀 Похоже, бот не привязан к учётной записи HMB Schools" +
		"\n\n" +
		"_Но если нужно привязать — оставляю эту кнопку здесь_ 👇"

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Войти",
					CallbackData: "/start",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructAttachedAndSend(ctx context.Context, chatIDs []int64, schoolName string) error {
	text := fmt.Sprintf("✅ Готово\\!\n\n"+
		"🔗 Привязали бота к *%s*\\. Можно начинать пользоваться\\.", schoolName)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Узнать дату следующей тренировки",
					CallbackData: "/next_training",
				},
			},
			{
				{
					Text:         "Узнать дату окончания абонемента",
					CallbackData: "/subscriptions",
				},
			},
			{
				{
					Text:         "Актуальное расписание",
					CallbackData: "/schedule",
				},
			},
		},
	}

	err := c.sendMessageWithButtonsToChatIDs(ctx, text, kb, chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) ConstructScheduleAndSend(ctx context.Context, chatIDs []int64, fioFull string, schoolName string, price int, desc string, schedule entity.Schedule) error {
	text := fmt.Sprintf(
		`🥊🥊🥊%s, актуальное расписание занятий в '%s':
`, fioFull, schoolName)
	for _, v := range schedule {
		text += fmt.Sprintf("\n🕗%s:%s, c %s по %s\n", v.Day, v.Description, v.Time.Open, v.Time.Closed)
	}
	text += fmt.Sprintf("\nцена месячного абонемента: 🔥%v р.\n\nДополнительное описание:\n%s", price, desc)
	err := c.sendMessageToChatIDs(ctx, bot.EscapeMarkdown(text), chatIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *constructor) sendMessageWithButtonsToChatIDs(ctx context.Context, text string, kb models.ReplyMarkup, chatIDs []int64) error {
	if len(chatIDs) == 0 {
		return nil
	}

	for _, chatID := range chatIDs {
		c.sender.SendMessageWithButtons(ctx, text, fmt.Sprint(chatID), kb)
	}

	return nil
}

func (c *constructor) sendPollToChatIDs(ctx context.Context, text string, opts []string, chatIDs []int64) error {
	if len(chatIDs) == 0 {
		return nil
	}

	for _, chatID := range chatIDs {
		c.sender.SendPoll(ctx, text, opts, fmt.Sprint(chatID))
	}

	return nil
}

func (c *constructor) sendMessageToChatIDs(ctx context.Context, text string, chatIDs []int64) error {
	if len(chatIDs) == 0 {
		return nil
	}

	for _, chatID := range chatIDs {
		c.sender.SendMessage(ctx, text, fmt.Sprint(chatID))
	}

	return nil
}

func getHelpBasicText() string {
	return "ℹ️ Это – бот HMB Schools\\. Через него мы рассылаем важные уведомления\\." +
		"\n\n" +
		"👉 Какие команды поддерживает этот бот:" +
		"\n\n" +
		"/start – подключиться к единой системе школ ИСБ России;\n" +
		"/register – самостоятельная регистрация ;\n" +
		"/info – показать, к какой учётной записи подключен бот;\n" +
		"/detach – отвязать бота от номера телефона;\n" +
		"/help – показать справку\\;"
}

func getHelpUserText() string {
	return `

	/subscriptions – получить актуальную дату окончания абонемента;

	/schedule – получить актуальное расписание зала вместе с ценой;

	/next_training – получить дату следующей тренировки;
`
}
func getHelpMasterText() string {
	return `

Мастер-команды:

	/subscription_list - получить список учеников с их актуальными абонементами;

	/hard_invite {username} {phone} {firstname} {lastname} – жесткий инвайт человека в вашу секцию;

	>>>Пример: /hard_invite - 71234567890 Иван Иванов


	/set_subscriptions {phone} {days} {*price} – Добавить абонемент;

	>>>Пример: /set_subscriptions 71234567890 1 500


	/set_schedule {day:openTime:closeTime:description;day:open:close:description;...} – Добавить расписание занятий;

	>>>Пример: /set_schedule пн:10-00:12-00:Борьба;ср:10-00:12-00:Кроссфит и прочие крутые штуки;пт:20-00:22-00: железо


	//set_desc_and_price_schedule {price} {description} – Изменить цену и описание занятий в месяц;
	
	>>>Пример: /set_desc_and_price_schedule 5000; у нас самый лучший тренер на планете земля Иван Иваныч


	/set_master {phone} – дать мастер-права для пользователя;
	
	>>>Пример: /set_master 71234567890 12


	/create_ticket {username} – создать токен для приглашения;(работает только с username)
	
	>>>Пример: /create_ticket  megaUser


P.S. username - как в телеге,если его нет - ставим прочерк('-') phone - формата 7ХХХХХХХХХХ
P.S.S где установлена '*', например *price - значит если вы его не введете, то значение будет по умолчанию в вашем зале.
`
}
