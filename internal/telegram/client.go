package telegram

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const CooldownAfter429 = 30 * time.Second

//go:generate mockery --name Client --case underscore --output mocks

type Client interface {
	SendMessage(ctx context.Context, text string, chatID string) (*models.Message, error)
	SendMessageWithButtons(ctx context.Context, text string, chatID string, kb models.ReplyMarkup) (*models.Message, error)
	SendPoll(ctx context.Context, text string, opts []string, chatID string) (*models.Message, error)
	SendPhoto(ctx context.Context, file []byte, caption string, chatID string) (*models.Message, error)
	EditMessageCaption(ctx context.Context, caption string, chatID string, messageID int) (*models.Message, error)
	EditMessageText(ctx context.Context, text string, chatID string, messageID int) (*models.Message, error)

	AddHandler(handle bot.HandlerFunc)
	Start(ctx context.Context)
}

type client struct {
	mx         *sync.Mutex
	cfg        ClientConfig
	botClient  *bot.Bot
	lastSentAt time.Time
}

func NewClient(cfg ClientConfig) (Client, error) {
	b, err := bot.New(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	return &client{
		mx:         &sync.Mutex{},
		cfg:        cfg,
		botClient:  b,
		lastSentAt: time.Now(),
	}, nil
}

func (c *client) AddHandler(handle bot.HandlerFunc) {
	c.botClient.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, handle)
	c.botClient.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, handle)
}

func (c *client) Start(ctx context.Context) {
	c.botClient.Start(ctx)
}

func (c *client) SendPoll(ctx context.Context, text string, opts []string, chatID string) (*models.Message, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	now := time.Now()
	if now.Before(c.lastSentAt.Add(c.cfg.PauseBetweenMessages)) {
		time.Sleep(c.lastSentAt.Add(c.cfg.PauseBetweenMessages).Sub(now))
	}
	c.lastSentAt = time.Now()
	options := make([]models.InputPollOption, 0, len(opts))
	for _, v := range opts {
		options = append(options, models.InputPollOption{Text: v})
	}

	m, err := c.botClient.SendPoll(ctx, &bot.SendPollParams{
		ChatID:                chatID,
		Question:              text,
		Options:               options,
		AllowsMultipleAnswers: true,
	})

	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			c.lastSentAt = time.Now().Add(CooldownAfter429)
		}
		return nil, err
	}

	return m, nil
}

func (c *client) SendMessage(ctx context.Context, text string, chatID string) (*models.Message, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	now := time.Now()
	if now.Before(c.lastSentAt.Add(c.cfg.PauseBetweenMessages)) {
		time.Sleep(c.lastSentAt.Add(c.cfg.PauseBetweenMessages).Sub(now))
	}
	c.lastSentAt = time.Now()

	m, err := c.botClient.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			c.lastSentAt = time.Now().Add(CooldownAfter429)
		}
		return nil, err
	}

	return m, nil
}

func (c *client) SendMessageWithButtons(ctx context.Context,
	text string, chatID string, kb models.ReplyMarkup) (*models.Message, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	now := time.Now()
	if now.Before(c.lastSentAt.Add(c.cfg.PauseBetweenMessages)) {
		time.Sleep(c.lastSentAt.Add(c.cfg.PauseBetweenMessages).Sub(now))
	}
	c.lastSentAt = time.Now()

	m, err := c.botClient.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: kb,
	})

	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			c.lastSentAt = time.Now().Add(CooldownAfter429)
		}
		return nil, err
	}

	return m, nil
}

func (c *client) SendPhoto(ctx context.Context,
	file []byte, caption string, chatID string) (*models.Message, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	now := time.Now()
	if now.Before(c.lastSentAt.Add(c.cfg.PauseBetweenMessages)) {
		time.Sleep(c.lastSentAt.Add(c.cfg.PauseBetweenMessages).Sub(now))
	}
	c.lastSentAt = time.Now()

	m, err := c.botClient.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: chatID,
		Photo: &models.InputFileUpload{
			Filename: "image.png",
			Data:     bytes.NewReader(file),
		},
		Caption:   caption,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			c.lastSentAt = time.Now().Add(CooldownAfter429)
		}
		return nil, err
	}

	return m, nil
}

//nolint:dupl // similar methods
func (c *client) EditMessageCaption(ctx context.Context,
	caption string, chatID string, messageID int) (*models.Message, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	now := time.Now()
	if now.Before(c.lastSentAt.Add(c.cfg.PauseBetweenMessages)) {
		time.Sleep(c.lastSentAt.Add(c.cfg.PauseBetweenMessages).Sub(now))
	}
	c.lastSentAt = time.Now()

	m, err := c.botClient.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
		ChatID:    chatID,
		MessageID: messageID,
		Caption:   caption,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			c.lastSentAt = time.Now().Add(CooldownAfter429)
		}
		return nil, err
	}

	return m, nil
}

//nolint:dupl // similar methods
func (c *client) EditMessageText(ctx context.Context,
	text string, chatID string, messageID int) (*models.Message, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	now := time.Now()
	if now.Before(c.lastSentAt.Add(c.cfg.PauseBetweenMessages)) {
		time.Sleep(c.lastSentAt.Add(c.cfg.PauseBetweenMessages).Sub(now))
	}
	c.lastSentAt = time.Now()

	m, err := c.botClient.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			c.lastSentAt = time.Now().Add(CooldownAfter429)
		}
		return nil, err
	}

	return m, nil
}
