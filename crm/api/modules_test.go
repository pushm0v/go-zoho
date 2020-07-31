package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpClient "github.com/pushm0v/go-zoho/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ZohoCrmAPiModulesSuite struct {
	suite.Suite
	url string
	api ApiModules
}

func TestZohoCrmApiModulesSuite(t *testing.T) {
	suite.Run(t, new(ZohoCrmAPiModulesSuite))
}

func (suite *ZohoCrmAPiModulesSuite) SetupTest() {
	sMock := suite.serverMock()
	suite.url = sMock.URL
	hClient := httpClient.NewHttpClient(sMock.Client())
	suite.api = NewApiModules(Option{
		ApiUrl:     suite.apiUrlMock,
		HttpClient: hClient,
	})
}

func (suite *ZohoCrmAPiModulesSuite) apiUrlMock(url string) string {
	return fmt.Sprintf("%s%s", suite.url, url)
}

func (suite *ZohoCrmAPiModulesSuite) serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(ZOHO_CRM_API_MODULES_URL, suite.modulesMock)

	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoCrmAPiModulesSuite) modulesMock(w http.ResponseWriter, r *http.Request) {
	var data = []byte(`{
		"modules": [{
			"id": "some-id",
			"api_name": "api-name",
			"module_name": "module-name"
		}]
	}`)

	_, _ = w.Write(data)
}

func (suite *ZohoCrmAPiModulesSuite) TestList() {
	modules, err := suite.api.List()
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), 1, len(modules), "Modules should not be empty")
	assert.Equal(suite.T(), "some-id", modules[0].ID, "ID not match")
	assert.Equal(suite.T(), "api-name", modules[0].ApiName, "API name not match")
	assert.Equal(suite.T(), "module-name", modules[0].ModuleName, "Module name not match")
}
