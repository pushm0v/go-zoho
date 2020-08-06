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
	assert.NotEmpty(suite.T(), suite.storage.Token.AccessToken(), "Access Token should not be empty")
	assert.NotEmpty(suite.T(), suite.storage.Token.RefreshToken(), "Refresh Token should not be empty")
}

func (suite *ZohoTokenStorageSuite) TestRefreshToken() {
	expectedRefreshToken := "some-refresh-token"
	suite.storage.Token.SaveToken("some-access-token", expectedRefreshToken, 60)
	assert.Equal(suite.T(), suite.storage.Token.RefreshToken(), expectedRefreshToken, "Refresh Token should be same")
}

func (suite *ZohoTokenStorageSuite) TestSaveAccessToken() {
	expectedToken := "some-access-token-new"
	suite.storage.Token.SaveAccessToken(expectedToken, 60)
	assert.Equal(suite.T(), suite.storage.Token.AccessToken(), expectedToken, "Access Token should be same")
}

func (suite *ZohoTokenStorageSuite) TestIsTokenExpired() {
	suite.storage.Token.SaveToken("some-access-token", "some-refresh-token", 3)
	assert.False(suite.T(), suite.storage.Token.IsTokenExpired(), "Token should not be expired")
}

func (suite *ZohoTokenStorageSuite) TestExpireTime() {
	expectedTime := 3
	suite.storage.Token.SaveToken("some-access-token", "some-refresh-token", expectedTime)

	assert.Equal(suite.T(), int(suite.storage.Token.ExpireTime()), expectedTime-1, "Expire time should be same")
}
