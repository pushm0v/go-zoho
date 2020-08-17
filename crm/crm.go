package crm

import (
	"fmt"

	"github.com/pushm0v/go-zoho/storage"

	"github.com/pushm0v/go-zoho/crm/api"
	"github.com/pushm0v/go-zoho/http"
)

type CrmVersion string

const (
	VersionV2 CrmVersion = "v2"
)

type ZohoCrmClient interface {
	Api(opts ...api.ApiOption) *api.CrmApi
}

type zohoCrmClient struct {
	Params ZohoCrmParams
	Option
}

type ZohoCrmParams struct {
	Version       CrmVersion
	ZGID          string
	CrmURL        string
	FileUploadURL string
}

type Option struct {
	HttpClient http.HttpClient
	Storage    *storage.Storage
}

func NewZohoCrmClient(params ZohoCrmParams, option Option) ZohoCrmClient {
	if params.Version == "" {
		params.Version = VersionV2
	}

	return &zohoCrmClient{
		Params: params,
		Option: option,
	}
}

func (z *zohoCrmClient) newApiOption() api.Option {
	z.HttpClient.WithAuthorizationFunc(z.getAccessToken)
	return api.Option{
		ApiUrl:        z.constructCrmUrl,
		FileUploadUrl: z.constructFileUploadUrl,
		HttpClient:    z.HttpClient,
		ApiParams:     z.getApiParams,
	}
}

func (z *zohoCrmClient) constructCrmUrl(url string, isBulk bool) string {
	if isBulk {
		return fmt.Sprintf("%s/bulk/%s%s", z.Params.CrmURL, z.Params.Version, url)
	}
	return fmt.Sprintf("%s/%s%s", z.Params.CrmURL, z.Params.Version, url)
}

func (z *zohoCrmClient) constructFileUploadUrl(url string) string {
	return fmt.Sprintf("%s/%s%s", z.Params.FileUploadURL, z.Params.Version, url)
}

func (z *zohoCrmClient) getApiParams(key string) interface{} {
	var apiParams = map[string]interface{}{
		"ZGID": z.Params.ZGID,
	}
	return apiParams[key]
}

func (z *zohoCrmClient) getAccessToken() string {
	return z.Storage.Token.AccessToken()
}

func (z *zohoCrmClient) Api(opts ...api.ApiOption) *api.CrmApi {
	c := &api.CrmApi{
		Option: z.newApiOption(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
