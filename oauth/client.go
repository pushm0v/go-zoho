package oauth

import (
	"encoding/json"
	"fmt"

	"github.com/pushm0v/go-zoho/http"
	"github.com/pushm0v/go-zoho/storage"
)

const (
	ZOHO_OAUTH_GRANT_URL         = "/oauth/v2/auth"
	ZOHO_OAUTH_TOKEN_URL         = "/oauth/v2/token"
	ZOHO_OAUTH_REFRESH_TOKEN_URL = "/oauth/v2/token"
	ZOHO_OAUTH_REVOKE_TOKEN_URL  = "/oauth/v2/token/revoke"
)

type ZohoAuthClient interface {
	GenerateToken(clientID, clientSecret, grantToken, iAMURL string) error
	RefreshToken() error
	OnSuccessTokenGeneration(f func(token OauthToken))
}

type zohoAuthClient struct {
	clientID      string
	clientSecret  string
	grantToken    string
	iAMURL        string
	httpClient    http.HttpClient
	storage       *storage.Storage
	onSuccessFunc func(token OauthToken)
}

func NewZohoAuthClient(httpClient http.HttpClient, storage *storage.Storage) ZohoAuthClient {
	return &zohoAuthClient{
		httpClient: httpClient,
		storage:    storage,
	}
}

func (z *zohoAuthClient) appendClientCredential(params map[string]interface{}) map[string]interface{} {
	params["client_id"] = z.clientID
	params["client_secret"] = z.clientSecret

	return params
}

func (z *zohoAuthClient) getURL(URL string) string {
	return fmt.Sprintf("%s%s", z.iAMURL, URL)
}

func (z *zohoAuthClient) GenerateToken(clientID, clientSecret, grantToken, iAMURL string) error {
	z.clientID = clientID
	z.clientSecret = clientSecret
	z.grantToken = grantToken
	z.iAMURL = iAMURL

	var params = map[string]interface{}{
		"grant_type": "authorization_code",
		"code":       grantToken,
	}
	params = z.appendClientCredential(params)

	resp, err := z.httpClient.Request("POST", z.getURL(ZOHO_OAUTH_TOKEN_URL), params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var token OauthToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return err
	}

	if token.Error != "" {
		return fmt.Errorf("Generate Token Error %v", token.Error)
	}

	z.saveToken(token)

	if z.onSuccessFunc != nil {
		z.onSuccessFunc(token)
	}
	return nil
}

func (z *zohoAuthClient) RefreshToken() error {
	var refreshToken = z.storage.Token.RefreshToken()

	if refreshToken == "" {
		return fmt.Errorf("Refresh Token is empty")
	}

	var params = map[string]interface{}{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	params = z.appendClientCredential(params)

	resp, err := z.httpClient.Request("POST", z.getURL(ZOHO_OAUTH_REFRESH_TOKEN_URL), params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var token OauthToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return err
	}

	if token.Error != "" {
		return fmt.Errorf("Refresh Token error %v", token.Error)
	}

	z.saveToken(token)

	if z.onSuccessFunc != nil {
		z.onSuccessFunc(token)
	}
	return nil
}

func (z *zohoAuthClient) OnSuccessTokenGeneration(f func(token OauthToken)) {
	z.onSuccessFunc = f
}

func (z *zohoAuthClient) saveToken(token OauthToken) {
	z.storage.Token.SaveToken(token.AccessToken, token.RefreshToken, token.ExpiresInSeconds)
	z.httpClient.WithAuthorization(token.AccessToken)
}
