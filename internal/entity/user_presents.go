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
