package storage

import (
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
	assert.NotNil(db)
}

func TestGetConnectionbyPublic(t *testing.T) {
	assert := assert.New(t)

	db, err := InitStorage(STORAGE_PATH)
	assert.NoError(err)

	publicToken := "public-token-2"

	db.NewConnectionWithToken("appid-1", "apptoken-1", "public-token-1", "<token>")
	db.NewConnectionWithToken("appid-1", "apptoken-2", publicToken, "<token>")

	conn := db.GetConnectionbyPublic(publicToken)
	assert.Equal(publicToken, conn.PublicToken)
}
