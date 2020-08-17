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
	GenerateToken() error
	RefreshToken() error
	OnSuccessTokenGeneration(f func(token OauthToken))
	TokenExpireTime() int
}

type zohoAuthClient struct {
	authParams    ZohoAuthParams
	httpClient    http.HttpClient
	storage       *storage.Storage
	onSuccessFunc func(token OauthToken)
}

type ZohoAuthParams struct {
	ClientID     string
	ClientSecret string
	GrantToken   string
	IamURL       string
}

func NewZohoAuthClient(authParams ZohoAuthParams, httpClient http.HttpClient, storage *storage.Storage) ZohoAuthClient {
	return &zohoAuthClient{
		authParams: authParams,
		httpClient: httpClient,
		storage:    storage,
	}
}

func (z *zohoAuthClient) appendClientCredential(params map[string]interface{}) map[string]interface{} {
	params["client_id"] = z.authParams.ClientID
	params["client_secret"] = z.authParams.ClientSecret

	return params
}

func (z *zohoAuthClient) getURL(URL string) string {
	return fmt.Sprintf("%s%s", z.authParams.IamURL, URL)
}

func (z *zohoAuthClient) isTokenExpired() bool {
	return z.storage.Token.IsTokenExpired()
}

func (z *zohoAuthClient) tokenFromStorage() OauthToken {
	accessToken := z.storage.Token.AccessToken()
	refreshToken := z.storage.Token.RefreshToken()
	expireTime := z.storage.Token.ExpireTime()

	return OauthToken{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresInSeconds: int(expireTime),
	}
}

func (z *zohoAuthClient) GenerateToken() error {

	var token OauthToken
	if z.isTokenExpired() {
		var params = map[string]interface{}{
			"grant_type": "authorization_code",
			"code":       z.authParams.GrantToken,
		}
		params = z.appendClientCredential(params)

		resp, err := z.httpClient.Post(z.getURL(ZOHO_OAUTH_TOKEN_URL), params)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&token)
		if err != nil {
			return err
		}

		if token.Error != "" {
			return fmt.Errorf("Generate Token Error %v", token.Error)
		}

		z.saveToken(token)
	} else {
		token = z.tokenFromStorage()
	}

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

	resp, err := z.httpClient.Post(z.getURL(ZOHO_OAUTH_REFRESH_TOKEN_URL), params)
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

	z.saveRefreshedToken(token)

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
}

func (z *zohoAuthClient) saveRefreshedToken(token OauthToken) {
	z.storage.Token.SaveAccessToken(token.AccessToken, token.ExpiresInSeconds)
}

func (z *zohoAuthClient) TokenExpireTime() int {
	return int(z.storage.Token.ExpireTime())
}
