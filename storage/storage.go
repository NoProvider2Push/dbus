package storage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitStorage(name string) *Storage {
	return &Storage{}
}

type Storage struct {
	db *gorm.DB
}

func (s Storage) NewApp(appID, token string) *AppToken {
	return &AppToken{
		db:          s.db,
		AppID:       appID,
		Token:       token,
		PublicToken: uuid.New().String(),
	}
}

type AppToken struct {
	db *gorm.DB

	AppID map[string]map[string]string
}

func (*AppToken) Add() {

}

func (*AppToken) Del() {

}
