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

func (s Storage) NewConnection(appID, token string, endpoint string) *Connection {

	return s.NewConnectionWithToken(appID, token, uuid.New().String(), endpoint)
}
func (s Storage) NewConnectionWithToken(appID, token string, publicToken, endpoint string) *Connection {
	existing := s.getFirst(Connection{AppID: appID, AppToken: token})
	if existing != nil {
		existing.PublicToken = publicToken
		existing.Endpoint = endpoint
		if err := s.db.Save(existing).Error; err != nil {
			return nil
		}
		return existing
	}

	// check if connection with this publicToken already exists
	// for pretend collision by with different app id and app token
	existing = s.getFirst(Connection{PublicToken: publicToken})
	if existing != nil {
		return nil
	}

	//create new if doesn't already exist
	c := Connection{
		AppID:       appID,
		AppToken:    token,
		PublicToken: publicToken,
		Endpoint:    endpoint,
	}
	result := s.db.Create(&c)
	if result.Error != nil {
		return nil
	}
	return &c
}

func (s Storage) DeleteConnection(token string) (*Connection, error) {
	c := Connection{AppToken: token}
	conn := s.getFirst(c)
	if conn == nil {
		return nil, errors.New("connection not found")
	}
	result := s.db.Delete(&c)
	return conn, result.Error
}

func (s Storage) GetConnectionbyPublic(publicToken string) *Connection {
	c := Connection{PublicToken: publicToken}
	return s.getFirst(c)
}

func (s Storage) GetUnequalEndpoint(latestEndpoint string) (ans []*Connection) {
	result := s.db.Find(&ans, "endpoint IS NOT ?", latestEndpoint)
	if result.Error != nil {
		return nil
	}
	return ans

}

func (s Storage) getFirst(c Connection) *Connection {
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
	Endpoint string
}
