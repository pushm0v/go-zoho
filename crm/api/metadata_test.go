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

type ZohoCrmAPiMetadataSuite struct {
	suite.Suite
	url string
	api ApiMetadata
}

func TestZohoCrmApiMetadataSuite(t *testing.T) {
	suite.Run(t, new(ZohoCrmAPiMetadataSuite))
}

func (suite *ZohoCrmAPiMetadataSuite) SetupTest() {
	sMock := suite.serverMock()
	suite.url = sMock.URL
	hClient := httpClient.NewHttpClient(sMock.Client())
	suite.api = NewApiMetadata(Option{
		ApiUrl:     suite.apiUrlMock,
		HttpClient: hClient,
	})
}

func (suite *ZohoCrmAPiMetadataSuite) apiUrlMock(url string) string {
	return fmt.Sprintf("%s%s", suite.url, url)
}

func (suite *ZohoCrmAPiMetadataSuite) serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(ZOHO_CRN_API_METADATA_FIELDS_URL, suite.fieldsMock)

	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoCrmAPiMetadataSuite) fieldsMock(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v", r.URL.String())
	keys, ok := r.URL.Query()["module"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'module' is missing")
		return
	}

	if keys[0] == "" {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`error`))
		return
	}

	var data = []byte(`{
		"fields": [{
			"id": "some-id",
			"api_name": "api-name",
			"field_label": "field-name"
		}]
	}`)

	_, _ = w.Write(data)
}

func (suite *ZohoCrmAPiMetadataSuite) TestListFields() {
	err, fields := suite.api.ListFields("test")
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), 1, len(fields), "Fields should not be empty")
}

func (suite *ZohoCrmAPiMetadataSuite) TestListFieldsNoModuleName() {
	err, fields := suite.api.ListFields("")
	assert.NotNil(suite.T(), err, "Error should be not nil")
	assert.Equal(suite.T(), 0, len(fields), "Fields should be empty")
}
