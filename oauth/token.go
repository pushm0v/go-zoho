package oauth

type OauthToken struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ApiDomain        string `json:"api_domain"`
	TokenType        string `json:"token_type"`
	ExpiresInSeconds int    `json:"expires_in"`
	Error            string `json:"error"`
}
