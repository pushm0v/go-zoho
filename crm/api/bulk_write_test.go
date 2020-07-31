package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpClient "github.com/pushm0v/go-zoho/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ZohoCrmAPIBulkWriteSuite struct {
	suite.Suite
	url string
	api ApiBulkWrite
}

func TestZohoCrmAPIBulkWriteSuite(t *testing.T) {
	suite.Run(t, new(ZohoCrmAPIBulkWriteSuite))
}

func (suite *ZohoCrmAPIBulkWriteSuite) SetupTest() {
	sMock := suite.serverMock()
	suite.url = sMock.URL
	hClient := httpClient.NewHttpClient(sMock.Client())
	suite.api = NewApiBulkWrite(Option{
		ApiUrl:     suite.apiUrlMock,
		FileUploadUrl:     suite.apiUrlMock,
		HttpClient: hClient,
	})
}

func (suite *ZohoCrmAPIBulkWriteSuite) apiUrlMock(url string) string {
	return fmt.Sprintf("%s%s", suite.url, url)
}

func (suite *ZohoCrmAPIBulkWriteSuite) serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(ZOHO_CRM_API_BULK_WRITE_UPLOAD_URL, suite.uploadZIPMock)

	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoCrmAPIBulkWriteSuite) uploadZIPMock(w http.ResponseWriter, r *http.Request) {
	var data = []byte(`{
		"status": "success",
		"code": "FILE_UPLOAD_SUCCESS",
		"message": "file uploaded.",
		"details": {
			"file_id": "123",
			"created_time": "2018-12-31T12:00:00-12:00"
		}
	}`)

	_, _ = w.Write(data)
}

func (suite *ZohoCrmAPIBulkWriteSuite) TestUploadZIP() {
	var fakeFile = strings.NewReader("fake, csv, data")
	fileID, err := suite.api.UploadZIP(fakeFile)
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), "123", fileID, "File ID not match")
}
