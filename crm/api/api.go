package api

import "github.com/pushm0v/go-zoho/http"

type Option struct {
	ApiUrl        ApiUrl
	FileUploadUrl FileUploadUrl
	HttpClient    http.HttpClient
	ApiParams     ApiParams
}

type ApiUrl func(string, bool) string
type FileUploadUrl func(string) string
type ApiOption func(*CrmApi)
type ApiParams func(string) interface{}

type CrmApi struct {
	Option          Option
	ApiBulkWrite    ApiBulkWrite
	ApiMetadata     ApiMetadata
	ApiModules      ApiModules
	ApiOrganization ApiOrganization
}

func WithApiBulkWrite() ApiOption {
	return func(a *CrmApi) {
		a.ApiBulkWrite = NewApiBulkWrite(a.Option)
	}
}
func WithApiModules() ApiOption {
	return func(a *CrmApi) {
		a.ApiModules = NewApiModules(a.Option)
	}
}

func WithApiMetadata() ApiOption {
	return func(a *CrmApi) {
		a.ApiMetadata = NewApiMetadata(a.Option)
	}
}

func WithApiOrganization() ApiOption {
	return func(a *CrmApi) {
		a.ApiOrganization = NewApiOrganization(a.Option)
	}
}
