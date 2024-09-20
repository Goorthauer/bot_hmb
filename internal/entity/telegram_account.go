package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type TelegramAccount struct {
	UserID    uuid.UUID `gorm:"column:user_id"`
	ChatID    int64     `gorm:"column:chat_id"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
