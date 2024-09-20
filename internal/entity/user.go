package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" gorm:"column:id"`
	IsMaster        bool       `json:"isMaster" gorm:"column:is_master"`
	IsDeleted       bool       `json:"is_deleted" gorm:"column:is_deleted"`
	IsActivated     bool       `json:"is_activated" gorm:"column:is_activated"`
	Username        string     `json:"username" gorm:"column:username"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
	DeletedBy       *uuid.UUID `json:"deleted_by" gorm:"column:deleted_by"`
	RegisteredAt    time.Time  `json:"registeredAt" gorm:"column:registered_at"`
	PDEncryptionKey *string    `json:"pd_encryption_key" gorm:"column:pd_encryption_key"`
	Phone           string     `json:"phone" gorm:"column:phone"`
	Region          string     `gorm:"<-:false;column:region"`

	SchoolID      uuid.UUID `gorm:"<-:false;column:school_id"`
	SchoolName    string    `gorm:"<-:false;column:school_name"`
	SchoolAddress string    `gorm:"<-:false;column:school_address"`

	PersonalData UserPersonalData `json:"personalData" gorm:"embedded"`
}

func (u *User) GetFullNameWithSchool() string {
	return fmt.Sprintf("%s (%s)", u.PersonalData.GetFullName(), u.SchoolName)
}
