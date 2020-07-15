package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pushm0v/go-zoho/storage"

	httpClient "github.com/pushm0v/go-zoho/http"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ZohoOauthClientSuite struct {
	suite.Suite
	httpClient httpClient.HttpClient
	storage    *storage.Storage
	params     TestZohoParams
}

type TestZohoParams struct {
	GrantToken   string
	ClientID     string
	ClientSecret string
	IAMURL       string
}

func TestZohoAuthClientSuite(t *testing.T) {
	suite.Run(t, new(ZohoOauthClientSuite))
}

func (suite *ZohoOauthClientSuite) SetupTest() {
	sMock := serverMock()
	suite.httpClient = httpClient.NewHttpClient(sMock.Client())
	suite.storage = storage.NewStorage()
	suite.params = TestZohoParams{
		GrantToken:   "some-token",
		ClientID:     "some",
		ClientSecret: "some",
		IAMURL:       sMock.URL,
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(ZOHO_OAUTH_TOKEN_URL, tokenMock)

	srv := httptest.NewServer(handler)

	return srv
}

func tokenMock(w http.ResponseWriter, r *http.Request) {
	var tokenResp = OauthToken{
		AccessToken:      "some-token",
		RefreshToken:     "some-refresh-token",
		ApiDomain:        "some-domain",
		ExpiresInSeconds: 1000,
		TokenType:        "some-token-type",
	}
	respByte, _ := json.Marshal(tokenResp)
	_, _ = w.Write(respByte)
}

func (suite *ZohoOauthClientSuite) TestGenerateToken() {
	client := NewZohoAuthClient(suite.httpClient, suite.storage)
	err := client.GenerateToken(suite.params.ClientID, suite.params.ClientSecret, suite.params.GrantToken, suite.params.IAMURL)
	assert.Nil(suite.T(), err, "Error should be nil")
}

func (suite *ZohoOauthClientSuite) TestRefreshToken() {
	client := NewZohoAuthClient(suite.httpClient, suite.storage)

	err := client.GenerateToken(suite.params.ClientID, suite.params.ClientSecret, suite.params.GrantToken, suite.params.IAMURL)
	assert.Nil(suite.T(), err, "Error should be nil")

	err = client.RefreshToken()
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), "some-refresh-token", suite.storage.Token.RefreshToken(), "Refresh Token is not same")
}

func (suite *ZohoOauthClientSuite) TestOnSuccessGenerateToken() {
	client := NewZohoAuthClient(suite.httpClient, suite.storage)
	expectedAccessToken := "some-token"
	success := func(t OauthToken) {
		assert.Equal(suite.T(), expectedAccessToken, t.AccessToken, "Access Token is not same")
	}
	client.OnSuccessTokenGeneration(success)
	err := client.GenerateToken(suite.params.ClientID, suite.params.ClientSecret, suite.params.GrantToken, suite.params.IAMURL)
	assert.Nil(suite.T(), err, "Error should be nil")
}
