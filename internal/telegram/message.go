package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	SendMessage = iota + 1
	SendMessageWithButtons
	SendPhoto
	SendPoll
	EditMessageCaption
	EditMessageText
	EditMessageTextWithButton
)

type Message struct {
	Operation      int
	ChatID         string
	Text           string
	Opts           []string
	FileArgName    []byte
	CaptionArgName string

	MessageID int

	Kb models.ReplyMarkup
}

type MessageQueue struct {
	MessageQueue chan Message
	Client       Client
	Constructor  Constructor
}

func NewMessageQueue(client Client) *MessageQueue {
	queue := &MessageQueue{
		MessageQueue: make(chan Message, 100),
		Client:       client,
	}
	go queue.QueueWorker()
	return queue
}

func (m *MessageQueue) QueueWorker() {
	for msg := range m.MessageQueue {
		var err error
		ctx := context.Background()
		switch msg.Operation {
		case SendMessage:
			_, err = m.Client.SendMessage(ctx, msg.Text, msg.ChatID)
		case SendPoll:
			_, err = m.Client.SendPoll(ctx, msg.Text, msg.Opts, msg.ChatID)
		case SendMessageWithButtons:
			_, err = m.Client.SendMessageWithButtons(ctx, msg.Text, msg.ChatID, msg.Kb)
		case SendPhoto:
			_, err = m.Client.SendPhoto(ctx, msg.FileArgName, msg.CaptionArgName, msg.ChatID)
		case EditMessageCaption:
			_, err = m.Client.EditMessageCaption(ctx, msg.CaptionArgName, msg.ChatID, msg.MessageID)
		case EditMessageText:
			_, err = m.Client.EditMessageText(ctx, msg.Text, msg.ChatID, msg.MessageID)
		case EditMessageTextWithButton:
			_, err = m.Client.EditMessageText(ctx, msg.Text, msg.ChatID, msg.MessageID)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}
