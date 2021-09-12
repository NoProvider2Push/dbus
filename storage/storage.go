package storage

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"unifiedpush.org/go/np2p_dbus/utils"
)

func InitStorage(name string) *Storage {
	db, err := gorm.Open(sqlite.Open(utils.StoragePath(name+".db")), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Connection{})
	return &Storage{db: db}
}

type Storage struct {
	db *gorm.DB
}

func (s Storage) NewConnection(appID, token string, endpoint string) *Connection {
	existing := s.getFirst(Connection{AppID: appID, AppToken: token})
	if existing != nil {
		return existing
	}

	//create new if doesn't already exist
	c := Connection{
		AppID:       appID,
		AppToken:    token,
		PublicToken: uuid.New().String(),
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
	result := s.db.First(&c)
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
