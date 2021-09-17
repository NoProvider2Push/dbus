package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const STORAGE_PATH = "/tmp/database.db"

func TestInit(t *testing.T) {
	assert := assert.New(t)

	db, err := InitStorage("/dev/notexistingfolder/database.db")
	assert.Error(err)
	assert.Nil(db)

	db, err = InitStorage(STORAGE_PATH)
	assert.NoError(err)
	defer os.Remove(STORAGE_PATH)
	assert.NotNil(db)
}

func TestNewConnectionWithGeneratedToken(t *testing.T) {
	assert := assert.New(t)

	db, err := InitStorage(STORAGE_PATH)
	assert.NoError(err)
	defer os.Remove(STORAGE_PATH)

	appID := "app-1"
	appToken := "apptoken-2"
	endpointSettings := "<endpoint>"

	conn := db.NewConnection(appID, appToken, endpointSettings)
	assert.NotNil(conn)
	// be sure that PublicToken is no given value
	assert.NotEqual("", conn.PublicToken)
	assert.NotEqual(appID, conn.PublicToken)
	assert.NotEqual(appToken, conn.PublicToken)
	assert.NotEqual(endpointSettings, conn.PublicToken)

	// that everythink else is given
	assert.Equal(appID, conn.AppID)
	assert.Equal(appToken, conn.AppToken)
	assert.Equal(endpointSettings, conn.Settings)
}

func TestNewConnectionUpdateSettings(t *testing.T) {
	assert := assert.New(t)

	db, err := InitStorage(STORAGE_PATH)
	assert.NoError(err)
	defer os.Remove(STORAGE_PATH)

	appID := "app-1"
	appToken := "apptoken-2"
	oldSettings := "endpoint-1"
	newSettings := "endpoint-2"

	// create connection
	conn := db.NewConnection(appID, appToken, oldSettings)
	assert.NotNil(conn)
	assert.Equal(oldSettings, conn.Settings)

	// save new settings on connection
	conn = db.NewConnection(appID, appToken, newSettings)
	assert.NotNil(conn)
	assert.Equal(newSettings, conn.Settings)
}

func TestGetConnectionbyPublic(t *testing.T) {
	assert := assert.New(t)

	db, err := InitStorage(STORAGE_PATH)
	assert.NoError(err)
	defer os.Remove(STORAGE_PATH)

	// create multiple connection
	db.NewConnection("appid-1", "apptoken-1", "<endpoint>")
	conn := db.NewConnection("appid-1", "apptoken-2", "<endpoint>")
	publicToken := conn.PublicToken
	db.NewConnection("appid-1", "apptoken-3", "<endpoint>")

	// find correct connection by public token
	conn = db.GetConnectionbyPublic(publicToken)
	assert.Equal(publicToken, conn.PublicToken)
}

func TestDeleteConnection(t *testing.T) {
	assert := assert.New(t)

	db, err := InitStorage(STORAGE_PATH)
	assert.NoError(err)
	defer os.Remove(STORAGE_PATH)

	appToken := "apptoken-2"

	// create multiple connection
	db.NewConnection("appid-1", "apptoken-1", "<endpoint>")
	db.NewConnection("appid-1", appToken, "<endpoint>")
	db.NewConnection("appid-1", "apptoken-3", "<endpoint>")

	// find correct connection by app token to delete
	conn, err := db.DeleteConnection(appToken)
	assert.NoError(err)
	assert.Equal(appToken, conn.AppToken)

	// unable to delete connection does not exists anymore
	conn, err = db.DeleteConnection(appToken)
	assert.Error(err)
}
