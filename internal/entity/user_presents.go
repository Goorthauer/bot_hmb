package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Presents struct {
	ID        uuid.UUID `json:"id" gorm:"column:id"`
	UserID    uuid.UUID `json:"user_id" gorm:"column:user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type PresentsSubscription struct {
	UserID        uuid.UUID `json:"user_id" gorm:"column:user_id"`
	CountTraining int       `json:"count_training" gorm:"column:count_training"`
	LostDays      int       `json:"subscription_days" gorm:"column:subscription_days"`
	DeadlineAt    time.Time `json:"deadline_at" gorm:"column:deadline_at"`
}
