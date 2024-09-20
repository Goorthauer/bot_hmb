package internal

import (
	"bot_hmb/internal/db"
	"bot_hmb/internal/dispatcher"
	"bot_hmb/internal/telegram"
	"bot_hmb/internal/usecase"
	"context"
	"log"
	"time"

	"bot_hmb/config"
)

type Manager struct {
	Queue *telegram.MessageQueue
}

func NewBotManager() *Manager {
	return &Manager{}
}

func (m *Manager) JoinBot() {
	conf := config.New()
	telegramCfg := telegram.ClientConfig{
		BotToken:             conf.Token,
		PauseBetweenMessages: 100 * time.Millisecond,
	}
	ctx := context.Background()
	b, err := telegram.NewClient(telegramCfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	conn, err := db.Connection(conf.Db)
	if err != nil {
		log.Fatal(err)
		return
	}
	redisClient, err := db.NewRedis(ctx, conf.RedisAddr)
	if err != nil {
		log.Fatal("redis not connected %w", err)
	}
	m.Queue = telegram.NewMessageQueue(b)
	wrapper := telegram.NewTelegramWrapper(m.Queue)
	constructor := telegram.NewConstructor(conf.Debug, wrapper)

	uc := usecase.New(conf.UserEncryptKey,
		constructor,
		wrapper,
		usecase.Config{
			TelegramBotURL:     conf.TelegramURL,
			MasterUserNickname: conf.MasterUserNickname},
		conn,
		redisClient)
	dp := dispatcher.NewDispatcher(uc)
	dp.SetHandlers()
	m.Queue.Client.AddHandler(dp.Dispatch)
	m.Queue.Client.Start(ctx)
}
