package api

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	ZOHO_CRM_API_BULK_WRITE_UPLOAD_URL     = "/upload"
	ZOHO_CRM_API_BULK_WRITE_CREATE_JOB_URL = "/write"
)

type JobOperationType string
type JobCallback struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}
type JobResource struct {
	Type        string `json:"type"`
	Module      string `json:"module"`
	FileID      string `json:"file_id"`
	IgnoreEmpty bool   `json:"ignore_empty"`
}

const (
	BulkWriteJobInsert JobOperationType = "insert"
	BulkWriteJobUpdate                  = "update"
	BulkWriteJobUpsert                  = "upsert"
)

type ApiBulkWriteJobParams struct {
	Operation JobOperationType `json:"operation"`
	Callback  JobCallback      `json:"callback"`
	Resource  JobResource      `json:"resource"`
}

type ApiBulkWrite interface {
	CreateJob(params ApiBulkWriteJobParams) (jobID string, err error)
	UploadZIP(handler io.Reader) (fileID string, err error)
}

type apiBulkWrite struct {
	option Option
}

type responseUploadZIP struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details struct {
		FileID string `json:"file_id"`
	} `json:"details"`
}

type responseBulkWriteCreate struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details struct {
		JobID string `json:"id"`
	} `json:"details"`
}

func NewApiBulkWrite(option Option) ApiBulkWrite {
	return &apiBulkWrite{
		option: option,
	}
}

func (bw *apiBulkWrite) CreateJob(params ApiBulkWriteJobParams) (jobID string, err error) {
	resp, err := bw.option.HttpClient.PostJson(bw.option.ApiUrl(ZOHO_CRM_API_BULK_WRITE_CREATE_JOB_URL), params)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var respWrite = new(responseBulkWriteCreate)

	err = json.NewDecoder(resp.Body).Decode(&respWrite)
	if err != nil {
		return
	}

	return respWrite.Details.JobID, nil
}

func (bw *apiBulkWrite) UploadZIP(handler io.Reader) (fileID string, err error) {
	var params = map[string]interface{}{}
	var headers = map[string]interface{}{
		"feature":   "bulk-write",
		"X-CRM-ORG": bw.option.ApiParams("ZGID"),
	}

	resp, err := bw.option.HttpClient.UploadZIP(bw.option.FileUploadUrl(ZOHO_CRM_API_BULK_WRITE_UPLOAD_URL), params, headers, handler)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var respUpload = new(responseUploadZIP)

	err = json.NewDecoder(resp.Body).Decode(&respUpload)
	if err != nil {
		return
	}

	if respUpload.Details.FileID == "" {
		return "", fmt.Errorf("%+v", respUpload)
	}

	return respUpload.Details.FileID, err
}
