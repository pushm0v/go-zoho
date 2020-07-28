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
	CrmURL        string
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
		ApiUrl: z.constructApiUrl,
		HttpClient: z.HttpClient,
	}
}

func (z *zohoCrmClient) constructApiUrl(url string) string {
	return fmt.Sprintf("%s%s", z.Params.CrmURL, url)
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
