package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ZohoHttpClientSuite struct {
	suite.Suite
	httpClient *http.Client
	httpServer *httptest.Server
}

func TestZohoHttpClientSuite(t *testing.T) {
	suite.Run(t, new(ZohoHttpClientSuite))
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoHttpClientSuite) SetupTest() {
	suite.httpClient = new(http.Client)
	suite.httpServer = serverMock()
}

func (suite *ZohoHttpClientSuite) TestPostRequest() {
	cHttp := NewHttpClient(suite.httpClient)
	expectedMethod := "POST"
	expectedURL := suite.httpServer.URL
	params := map[string]interface{}{
		"some-key": 1,
	}
	resp, err := cHttp.Post(expectedURL, params)
	assert.NoError(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), expectedMethod, resp.Request.Method)
	assert.Equal(suite.T(), expectedURL, resp.Request.URL.String())
}

func (suite *ZohoHttpClientSuite) TestPostJsonRequest() {
	cHttp := NewHttpClient(suite.httpClient)
	expectedMethod := "POST"
	expectedContentType := "application/json"
	expectedURL := suite.httpServer.URL
	params := map[string]interface{}{
		"some-key": 1,
	}
	resp, err := cHttp.PostJson(expectedURL, params)
	assert.NoError(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), expectedMethod, resp.Request.Method)
	assert.Equal(suite.T(), expectedURL, resp.Request.URL.String())
	assert.Equal(suite.T(), expectedContentType, resp.Request.Header["Content-Type"][0])
}

func (suite *ZohoHttpClientSuite) TestGetRequest() {
	cHttp := NewHttpClient(suite.httpClient)
	expectedMethod := "GET"
	expectedURL := fmt.Sprintf("%s?some-key=1", suite.httpServer.URL)
	params := map[string]interface{}{
		"some-key": "1",
	}
	resp, err := cHttp.Get(suite.httpServer.URL, params)
	assert.NoError(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), expectedMethod, resp.Request.Method)
	assert.Equal(suite.T(), expectedURL, resp.Request.URL.String())
	assert.Equal(suite.T(), "1", resp.Request.URL.Query().Get("some-key"))
}

func (suite *ZohoHttpClientSuite) TestUploadZIPRequest() {
	cHttp := NewHttpClient(suite.httpClient)
	expectedMethod := "POST"
	expectedURL := suite.httpServer.URL
	params := map[string]interface{}{
		"some-key": 1,
	}
	headers := map[string]interface{}{
		"some-headers": "some-value",
	}
	var fakeFile = strings.NewReader("fake, csv, data")
	resp, err := cHttp.UploadZIP(expectedURL, params, headers, fakeFile)
	assert.NoError(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), expectedMethod, resp.Request.Method)
	assert.Equal(suite.T(), expectedURL, resp.Request.URL.String())
	assert.Equal(suite.T(), "some-value", resp.Request.Header.Get("some-headers"))
}

func (suite *ZohoHttpClientSuite) TestWithAuthorizationFunc() {
	cHttp := NewHttpClient(suite.httpClient)
	expectedAccessToken := "some-token"
	var f = func() string {
		return expectedAccessToken
	}
	cHttp.WithAuthorizationFunc(f)
	expectedMethod := "GET"
	expectedURL := suite.httpServer.URL
	params := map[string]interface{}{}
	resp, err := cHttp.Get(suite.httpServer.URL, params)
	assert.NoError(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), expectedMethod, resp.Request.Method)
	assert.Equal(suite.T(), expectedURL, resp.Request.URL.String())
	assert.Equal(suite.T(), fmt.Sprintf("Zoho-oauthtoken %s", expectedAccessToken), resp.Request.Header["Authorization"][0])
}
