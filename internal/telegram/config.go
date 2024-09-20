package telegram

import (
	"time"
)

type ClientConfig struct {
	BotToken             string
	PauseBetweenMessages time.Duration
}
