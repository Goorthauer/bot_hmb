package telegram

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-telegram/bot/models"
)

type Wrapper interface {
	SendMessage(ctx context.Context, text string, chatID string)
	SendMessageWithButtons(ctx context.Context, text string, chatID string, kb models.ReplyMarkup)
	SendPhoto(ctx context.Context, file []byte, caption string, chatID string)
	SendPoll(ctx context.Context, text string, opts []string, chatID string)
	EditMessageCaption(ctx context.Context, caption string, chatID string, messageID int)
	EditMessageText(ctx context.Context, text string, chatID string, messageID int)
	EditMessageTextWithButtons(ctx context.Context, text string, chatID string, messageID int, kb models.ReplyMarkup)
}

type telegramWrapper struct {
	Queue *MessageQueue
}

func NewTelegramWrapper(queue *MessageQueue) Wrapper {
	return &telegramWrapper{
		Queue: queue,
	}
}

func (w *telegramWrapper) SendMessage(_ context.Context, text string, chatID string) {
	w.Queue.MessageQueue <- Message{
		Operation: SendMessage,
		ChatID:    chatID,
		Text:      text,
	}
}
func (w *telegramWrapper) SendPoll(_ context.Context, text string, opts []string, chatID string) {
	w.Queue.MessageQueue <- Message{
		Operation: SendPoll,
		ChatID:    chatID,
		Text:      text,
		Opts:      opts,
	}
}

func (w *telegramWrapper) SendMessageWithButtons(_ context.Context,
	text string, chatID string, keyboard models.ReplyMarkup) {
	kb, err := json.Marshal(keyboard)
	if err != nil {
		fmt.Println("kb marshall err: %w", err)
		return
	}
	w.Queue.MessageQueue <- Message{
		Operation: SendMessageWithButtons,
		ChatID:    chatID,
		Text:      text,
		Kb:        string(kb),
	}
}

func (w *telegramWrapper) SendPhoto(_ context.Context, file []byte, caption string, chatID string) {
	w.Queue.MessageQueue <- Message{
		Operation:      SendPhoto,
		ChatID:         chatID,
		FileArgName:    file,
		CaptionArgName: caption,
	}
}

func (w *telegramWrapper) EditMessageCaption(_ context.Context, caption string, chatID string, messageID int) {
	w.Queue.MessageQueue <- Message{
		Operation:      EditMessageCaption,
		ChatID:         chatID,
		CaptionArgName: caption,
		MessageID:      messageID,
	}
}

func (w *telegramWrapper) EditMessageText(_ context.Context, text string, chatID string, messageID int) {
	w.Queue.MessageQueue <- Message{
		Operation: EditMessageText,
		ChatID:    chatID,
		Text:      text,
		MessageID: messageID,
	}
}

func (w *telegramWrapper) EditMessageTextWithButtons(_ context.Context, text string, chatID string, messageID int, keyboard models.ReplyMarkup) {
	kb, err := json.Marshal(keyboard)
	if err != nil {
		fmt.Println("kb marshall err: %w", err)
		return
	}

	w.Queue.MessageQueue <- Message{
		Operation: EditMessageTextWithButton,
		ChatID:    chatID,
		Text:      text,
		MessageID: messageID,
		Kb:        kb,
	}
}
