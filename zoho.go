package zoho

import (
	"net/http"

	"github.com/pushm0v/go-zoho/crm"

	http2 "github.com/pushm0v/go-zoho/http"
	"github.com/pushm0v/go-zoho/oauth"
	storage2 "github.com/pushm0v/go-zoho/storage"
)

type Zoho interface {
	Authenticate() error
	Reauthenticate() error
	Crm() crm.ZohoCrmClient
}

type zoho struct {
	params ZohoParams
	oauth  oauth.ZohoAuthClient
	crm    crm.ZohoCrmClient
}

type ZohoParams struct {
	GrantToken   string
	ClientID     string
	ClientSecret string
	IamURL       string
	CrmURL       string
}

func NewZoho(params ZohoParams, httpClient *http.Client) Zoho {
	var authParams = oauth.ZohoAuthParams{
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		GrantToken:   params.GrantToken,
		IamURL:       params.IamURL,
	}
	var crmParams = crm.ZohoCrmParams{
		CrmURL: params.CrmURL,
	}
	var hClient http2.HttpClient
	if httpClient == nil {
		hClient = defaultHttpClient()
	} else {
		hClient = http2.NewHttpClient(httpClient)
	}
	var storage = defaultStorage()

	return &zoho{
		oauth: newOauthClient(authParams, hClient, storage),
		crm:   newCrmClient(crmParams, hClient),
	}
}

func defaultStorage() *storage2.Storage {
	return storage2.NewStorage()
}

func defaultHttpClient() http2.HttpClient {
	return http2.NewHttpClient(new(http.Client))
}

func newOauthClient(authParams oauth.ZohoAuthParams, httpClient http2.HttpClient, storage *storage2.Storage) oauth.ZohoAuthClient {
	return oauth.NewZohoAuthClient(authParams, httpClient, storage)
}

func newCrmClient(crmParams crm.ZohoCrmParams, httpClient http2.HttpClient) crm.ZohoCrmClient {
	return crm.NewZohoCrmClient(crmParams, crm.Option{HttpClient: httpClient})
}

func (z *zoho) Authenticate() error {
	return z.oauth.GenerateToken()
}

func (z *zoho) Reauthenticate() error {
	return z.oauth.RefreshToken()
}

func (z *zoho) Crm() crm.ZohoCrmClient {
	return z.crm
}
