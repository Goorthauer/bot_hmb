package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type TelegramAuthTicket struct {
	Token     string     `gorm:"column:token"`
	UserID    uuid.UUID  `gorm:"column:user_id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	IsSpent   bool       `gorm:"column:is_spent"`
	SpentAt   *time.Time `gorm:"column:spent_at"`
	IsBlocked bool       `gorm:"column:is_blocked"`
}
