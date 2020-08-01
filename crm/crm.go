package crm

import (
	"fmt"

	"github.com/pushm0v/go-zoho/crm/api"
	"github.com/pushm0v/go-zoho/http"
)

type ZohoCrmClient interface {
	Api(opts ...api.ApiOption) *api.CrmApi
}

type zohoCrmClient struct {
	Params ZohoCrmParams
	Option
}

type ZohoCrmParams struct {
	ZGID          string
	CrmURL        string
	FileUploadURL string
}

type Option struct {
	HttpClient http.HttpClient
}

func NewZohoCrmClient(params ZohoCrmParams, option Option) ZohoCrmClient {
	return &zohoCrmClient{
		Params: params,
		Option: option,
	}
}

func (z *zohoCrmClient) newApiOption() api.Option {
	return api.Option{
		ApiUrl:        z.constructApiUrl,
		FileUploadUrl: z.constructFileUPloadUrl,
		HttpClient:    z.HttpClient,
		ApiParams:     z.getApiParams,
	}
}

func (z *zohoCrmClient) constructApiUrl(url string) string {
	return fmt.Sprintf("%s%s", z.Params.CrmURL, url)
}

func (z *zohoCrmClient) constructFileUPloadUrl(url string) string {
	return fmt.Sprintf("%s%s", z.Params.FileUploadURL, url)
}

func (z *zohoCrmClient) getApiParams(key string) interface{} {
	var apiParams = map[string]interface{}{
		"ZGID": z.Params.ZGID,
	}
	return apiParams[key]
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
