package http

import (
	"net/http"
	"net/http/httptest"
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

func (suite *ZohoHttpClientSuite) TestRequest() {
	cHttp := NewHttpClient(suite.httpClient)
	expectedMethod := "POST"
	expectedURL := suite.httpServer.URL
	params := map[string]interface{}{
		"some-key": 1,
	}
	resp, err := cHttp.Request(expectedMethod, expectedURL, params)
	assert.NoError(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), expectedMethod, resp.Request.Method)
	assert.Equal(suite.T(), expectedURL, resp.Request.URL.String())
}
