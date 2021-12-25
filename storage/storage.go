package storage

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitStorage(filepath string) (*Storage, error) {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger: logger.Default.LogMode(0),
	})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Connection{})
	return &Storage{db: db}, nil
}

type Storage struct {
	db *gorm.DB
}

func (s Storage) DB() *gorm.DB {
	return s.db
}

func (s Storage) NewConnection(appID, appToken string, settings string) *Connection {
	return s.NewConnectionFull(appID, appToken, uuid.New().String(), settings)
}

func (s Storage) NewConnectionFull(appID, appToken, pubToken, settings string) *Connection {
	existing := s.GetFirst(Connection{AppID: appID, AppToken: appToken})
	if existing != nil {
		existing.Settings = settings
		if err := s.db.Save(existing).Error; err != nil {
			return nil
		}
		return existing
	}

	//create new if doesn't already exist
	c := Connection{
		AppID:       appID,
		AppToken:    appToken,
		PublicToken: pubToken,
		Settings:    settings,
	}
	result := s.db.Create(&c)
	if result.Error != nil {
		return nil
	}
	return &c
}

func (s Storage) DeleteConnection(token string) (*Connection, error) {
	c := Connection{AppToken: token}
	conn := s.GetFirst(c)
	if conn == nil {
		return nil, errors.New("connection not found")
	}
	result := s.db.Delete(&c)
	return conn, result.Error
}

func (s Storage) GetConnectionbyPublic(publicToken string) *Connection {
	c := Connection{PublicToken: publicToken}
	return s.GetFirst(c)
}

func (s Storage) GetUnequalSettings(latestSettings string) (ans []*Connection) {
	result := s.db.Find(&ans, "settings IS NOT ?", latestSettings)
	if result.Error != nil {
		return nil
	}
	return ans

}

func (s Storage) GetFirst(c Connection) *Connection {
	result := s.db.Where(&c).First(&c)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &c
}

type Connection struct {
	AppID       string
	AppToken    string `gorm:"primaryKey"`
	PublicToken string `gorm:"unique"`
	// endpoint format, to keep track of settings changes
	Settings string
}
