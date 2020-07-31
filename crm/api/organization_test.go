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

type ZohoCrmAPiOrganizationSuite struct {
	suite.Suite
	url string
	api ApiOrganization
}

func TestZohoCrmApiOrganizationSuite(t *testing.T) {
	suite.Run(t, new(ZohoCrmAPiOrganizationSuite))
}

func (suite *ZohoCrmAPiOrganizationSuite) SetupTest() {
	sMock := suite.serverMock()
	suite.url = sMock.URL
	hClient := httpClient.NewHttpClient(sMock.Client())
	suite.api = NewApiOrganization(Option{
		ApiUrl:     suite.apiUrlMock,
		HttpClient: hClient,
	})
}

func (suite *ZohoCrmAPiOrganizationSuite) apiUrlMock(url string) string {
	return fmt.Sprintf("%s%s", suite.url, url)
}

func (suite *ZohoCrmAPiOrganizationSuite) serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(ZOHO_CRM_API_ORGANIZATION_URL, suite.organizationMock)

	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoCrmAPiOrganizationSuite) organizationMock(w http.ResponseWriter, r *http.Request) {
	var data = []byte(`{
		"org": [{
			"id": "some-id",
			"zgid": "123"
		}]
	}`)

	_, _ = w.Write(data)
}

func (suite *ZohoCrmAPiOrganizationSuite) TestDetails() {
	err, details := suite.api.Details()
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), 1, len(details), "Organization should not be empty")
	assert.Equal(suite.T(), "123", details[0].ZGID, "ZGID not match")
	assert.Equal(suite.T(), "some-id", details[0].ID, "ID not match")
}
