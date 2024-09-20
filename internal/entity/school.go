package entity

import (
	"github.com/gofrs/uuid"
)

type School struct {
	ID      uuid.UUID `json:"id" gorm:"column:id"`
	Region  int       `json:"region" gorm:"column:region"`
	Name    string    `json:"name" gorm:"column:name"`
	City    string    `json:"city" gorm:"column:city"`
	Address string    `json:"address" gorm:"column:address"`
	Contact string    `json:"contact" gorm:"column:contact"`
	VkLink  string    `json:"vk_link" gorm:"column:vk_link"`
}
