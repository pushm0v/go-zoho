package zoho

import (
	"net/http"

	http2 "github.com/pushm0v/go-zoho/http"
	"github.com/pushm0v/go-zoho/oauth"
	storage2 "github.com/pushm0v/go-zoho/storage"
)

type Zoho interface {
	Authenticate() error
	WithHttpClient(httpClient *http.Client)
}

type zoho struct {
	params ZohoParams
	oauth  oauth.ZohoAuthClient
}

type ZohoParams struct {
	GrantToken   string
	ClientID     string
	ClientSecret string
	IAMURL       string
}

func NewZoho(params ZohoParams) Zoho {
	return &zoho{
		params: params,
		oauth:  newOauthClient(nil),
	}
}

func defaultStorage() *storage2.Storage {
	return storage2.NewStorage()
}

func defaultHttpClient() http2.HttpClient {
	return http2.NewHttpClient(new(http.Client))
}

func newOauthClient(httpClient *http.Client) oauth.ZohoAuthClient {
	var hClient http2.HttpClient
	if httpClient == nil {
		hClient = defaultHttpClient()
	} else {
		hClient = http2.NewHttpClient(httpClient)
	}

	var storage = defaultStorage()
	return oauth.NewZohoAuthClient(hClient, storage)
}

func (z *zoho) Authenticate() error {
	return z.oauth.GenerateToken(z.params.ClientID, z.params.ClientSecret, z.params.GrantToken, z.params.IAMURL)
}

func (z *zoho) WithHttpClient(httpClient *http.Client) {
	z.oauth = newOauthClient(httpClient)
}
