package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Manager struct {
	Gorm *gorm.DB
}

func Connection(dsn string) (*Manager, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return &Manager{Gorm: db}, err
}
