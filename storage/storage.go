package storage

import (
	"github.com/google/uuid"
	"github.com/karmanyaahm/np2p_linux/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitStorage(name string) *Storage {
	db, err := gorm.Open(sqlite.Open(utils.StoragePath("np2p")), &gorm.Config{})
	if err != nil {
		log.Fataln("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Connection{})
	return &Storage{db: db}
}

type Storage struct {
	db *gorm.DB
}

func (s Storage) NewConnection(appID, token string) *Connection {
	c := Connection{
		db:          s.db,
		AppID:       appID,
		AppToken:    token,
		PublicToken: uuid.New().String(),
	}
	result := s.db.Create(&c)
	if result.Error != nil {
		return nil
	}
	return &c
}

func (s Storage) DeleteConnection(appID, token string) error {
	result := s.db.Delete(&Connection{AppID: appID, AppToken: token})
	return result.Error
}

func (s Storage) GetConnectionbyPublic(publicToken string) *Connection {
	c := Connection{}
	result := s.db.First(&c)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &c
}

type Connection struct {
	db *gorm.DB

	AppID       string `gorm:"primaryKey"`
	AppToken    string `gorm:"primaryKey"`
	PublicToken string `gorm:"unique"`
}
