package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
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
		ApiUrl:        suite.apiUrlMock,
		FileUploadUrl: suite.fileUploadUrlMock,
		HttpClient:    hClient,
		ApiParams:     suite.apiParamsMock,
	})
}

func (suite *ZohoCrmAPIBulkWriteSuite) apiParamsMock(key string) interface{} {
	return key
}

func (suite *ZohoCrmAPIBulkWriteSuite) apiUrlMock(url string, isBulk bool) string {
	return fmt.Sprintf("%s%s", suite.url, url)
}

func (suite *ZohoCrmAPIBulkWriteSuite) fileUploadUrlMock(url string) string {
	return fmt.Sprintf("%s%s", suite.url, url)
}

func (suite *ZohoCrmAPIBulkWriteSuite) serverMock() *httptest.Server {
	handler := mux.NewRouter()
	handler.HandleFunc(ZOHO_CRM_API_BULK_WRITE_UPLOAD_URL, suite.uploadZIPMock)
	handler.HandleFunc(ZOHO_CRM_API_BULK_WRITE_CREATE_JOB_URL, suite.createJobMock)
	handler.HandleFunc(fmt.Sprintf("%s/{jobID}", ZOHO_CRM_API_BULK_WRITE_CREATE_JOB_URL), suite.jobDetailsMock)

	srv := httptest.NewServer(handler)

	return srv
}

func (suite *ZohoCrmAPIBulkWriteSuite) createJobMock(w http.ResponseWriter, r *http.Request) {
	var data = []byte(`{
	  "status": "success",
	  "code": "SUCCESS",
	  "message": "success",
	  "details": {
		"id": "111111000000541958",
		"created_by": {
		  "id": "111111000000035795",
		  "name": "Patricia Boyle "
		}
	  }
	}`)

	_, _ = w.Write(data)
}

func (suite *ZohoCrmAPIBulkWriteSuite) jobDetailsMock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	jobID := vars["jobID"]
	var data = []byte(fmt.Sprintf(`{
	  "status": "COMPLETED",
	  "character_encoding": "UTF-8",
	  "resource": [
		{
		  "status": "COMPLETED",
		  "type": "data",
		  "module": "Deals",
		  "field_mappings": [
			{
			  "api_name": "Deal_Name",
			  "index": 1,
			  "format": null,
			  "find_by": null,
			  "default_value": null
			},
			{
			  "api_name": "Stage",
			  "index": 2,
			  "format": null,
			  "find_by": null,
			  "default_value": null
			}
		  ],
		  "file": {
			"status": "COMPLETED",
			"name": "Accounts.csv",
			"added_count": 0,
			"skipped_count": 100,
			"updated_count": 0,
			"total_count": 100
		  }
		}
	  ],
	  "id": "%s",
	  "result": {
		"download_url": "/v2/crm/org6196138/bulk-write/111111000002308051/111111000002308051.zip"
	  },
	  "created_by": {
		"id": "111111000000035795",
		"name": "Patricia Boyle"
	  },
	  "operation": "insert",
	  "created_time": "2019-01-30T02:18:15-12:00"
	}`, jobID))

	_, _ = w.Write(data)
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

func (suite *ZohoCrmAPIBulkWriteSuite) TestCerateJob() {
	var fakeParams = ApiBulkWriteJobParams{
		Operation: BulkWriteJobInsert,
		Callback: JobCallback{
			Method: "POST",
			Url:    "http://some-url",
		},
		Resource: []JobResource{{
			Type:        "some-type",
			Module:      "some-module",
			FileID:      "123",
			IgnoreEmpty: true,
		}},
	}
	JobID, err := suite.api.CreateJob(fakeParams)
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), "111111000000541958", JobID, "Job ID not match")
}

func (suite *ZohoCrmAPIBulkWriteSuite) TestJobDetails() {

	result, err := suite.api.JobDetails("123")
	assert.Nil(suite.T(), err, "Error should be nil")
	assert.NotEmpty(suite.T(), result, "Result should not be")
	assert.Equal(suite.T(), "123", result["id"], "JobID not match")
}
