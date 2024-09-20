package entity

import (
	"github.com/gofrs/uuid"
)

type UserSchool struct {
	UserID   uuid.UUID `json:"user_id" gorm:"column:user_id"`
	SchoolID uuid.UUID `json:"school_id" gorm:"column:school_id"`
}
