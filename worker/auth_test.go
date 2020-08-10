package worker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpClient "github.com/pushm0v/go-zoho/http"
	"github.com/stretchr/testify/assert"

	"github.com/pushm0v/go-zoho/cache"
	"github.com/pushm0v/go-zoho/oauth"
	storage2 "github.com/pushm0v/go-zoho/storage"
	"github.com/stretchr/testify/suite"
)

type ZohoAuthWorkerSuite struct {
	suite.Suite
	worker AuthWorker
	params AuthWorkerParams
}

func TestZohoAuthWorkerSuiteSuite(t *testing.T) {
	suite.Run(t, new(ZohoAuthWorkerSuite))
}

func (suite *ZohoAuthWorkerSuite) tokenMock(w http.ResponseWriter, r *http.Request) {
	var tokenResp = oauth.OauthToken{
		AccessToken:      "some-token",
		RefreshToken:     "some-refresh-token",
		ApiDomain:        "some-domain",
		ExpiresInSeconds: 5,
		TokenType:        "some-token-type",
	}
	respByte, _ := json.Marshal(tokenResp)
	_, _ = w.Write(respByte)
}

func (suite *ZohoAuthWorkerSuite) serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(oauth.ZOHO_OAUTH_TOKEN_URL, suite.tokenMock)

	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoAuthWorkerSuite) SetupTest() {
	suite.params = AuthWorkerParams{
		SecondsBeforeRefreshToken: 3,
	}
	sMock := suite.serverMock()
	hClient := httpClient.NewHttpClient(sMock.Client())
	oParams := oauth.ZohoAuthParams{
		GrantToken:   "some-token",
		ClientID:     "some",
		ClientSecret: "some",
		IamURL:       sMock.URL,
	}
	storage := storage2.NewStorage(cache.NewCache(cache.WithLocalCache()))
	client := oauth.NewZohoAuthClient(oParams, hClient, storage)

	suite.worker = NewAuthWorker(client, suite.params)
}

func (suite *ZohoAuthWorkerSuite) TestAuthWorkerStart() {
	onErrorFunc := func(err error) {
		assert.NotNil(suite.T(), err, "Error should be not nil")
	}
	suite.worker.OnError(onErrorFunc)
	suite.worker.Start()
}
