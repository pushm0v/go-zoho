package api

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	ZOHO_CRM_API_BULK_WRITE_UPLOAD_URL = "/upload"
)

type ApiBulkWrite interface {
	Write()
	UploadZIP(handler io.Reader) (fileID string, err error)
}

type apiBulkWrite struct {
	option Option
}

type responseUploadZIP struct {
	Status string `json:"status"`
	Code string `json:"code"`
	Message string `json:"message"`
	Details struct {
		FileID string `json:"file_id"`
		CreatedTime string `json:"created_time"`
	} `json:"details"`
}

func NewApiBulkWrite(option Option) ApiBulkWrite {
	return &apiBulkWrite{
		option: option,
	}
}

func (bw *apiBulkWrite) Write() {

}

func (bw *apiBulkWrite) UploadZIP(handler io.Reader) (fileID string, err error) {
	var params = map[string]interface{}{}
	var headers = map[string]interface{}{}

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
		return "", fmt.Errorf("%v", respUpload.Message)
	}

	return respUpload.Details.FileID, err
}
