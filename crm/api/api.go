package api

import "github.com/pushm0v/go-zoho/http"

type Option struct {
	ApiUrl     ApiUrl
	HttpClient http.HttpClient
}

type ApiUrl func(string) string
type ApiOption func(*CrmApi)

type CrmApi struct {
	Option       Option
	ApiBulkWrite ApiBulkWrite
	ApiMetadata  ApiMetadata
	ApiModules   ApiModules
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
