package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Subscription struct {
	ID         uuid.UUID `json:"id" gorm:"column:id"`
	UserID     uuid.UUID `json:"user_id" gorm:"column:user_id"`
	SchoolID   uuid.UUID `json:"school_id" gorm:"column:school_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	DeadlineAt time.Time `json:"deadline_at" gorm:"column:deadline_at"`
	Price      string    `json:"price" gorm:"column:price"`
	Days       int       `json:"days" gorm:"column:days"`
}
