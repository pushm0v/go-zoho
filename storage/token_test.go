package storage

import (
	"testing"

	"github.com/pushm0v/go-zoho/cache"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ZohoTokenStorageSuite struct {
	suite.Suite
	storage *Storage
}

func TestZohoTokenStorageSuite(t *testing.T) {
	suite.Run(t, new(ZohoTokenStorageSuite))
}

func (suite *ZohoTokenStorageSuite) SetupTest() {
	suite.storage = NewStorage(cache.NewCache(cache.WithLocalCache()))
}

func (suite *ZohoTokenStorageSuite) TestSave() {
	suite.storage.Token.SaveToken("some-access-token", "some-refresh-token", 60)
	assert.NotEmpty(suite.T(), suite.storage.Token.AccessToken(), "Token should not be empty")
}

func (suite *ZohoTokenStorageSuite) TestRefreshToken() {
	suite.storage.Token.SaveToken("some-access-token", "some-refresh-token", 60)
	assert.NotEmpty(suite.T(), suite.storage.Token.RefreshToken(), "Refresh Token should not be empty")
}
