package github

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

const (
	clientId           = "GH_CLIENT_ID"
	clientSecret       = "GH_CLIENT_SECRET"
	authorizeUrl       = "GH_AUTHORIZE_URL"
	authAccessTokenUrl = "GH_ACCESS_TOKEN_URL"
	authCallbackUrl    = "GH_AUTH_CALLBACK_URL"
	headerAcceptParam  = "Accept"
	headerAcceptValue  = "application/json"
	authCodeParam      = "code"
)

type ghCredentials struct {
	clientId     string
	clientSecret string
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

var credentials = getCredentials()
var httpClient = http.Client{}

func AuthorizationUrl() string {
	return fmt.Sprintf(viper.GetString(authorizeUrl), credentials.clientId)
}

func Authorize(req *http.Request) (OAuthAccessResponse, error) {
	code := req.URL.Query().Get(authCodeParam)
	accessTokenUrl := getAccessTokenUrl(code)
	accessTokenRequest, _ := http.NewRequest(http.MethodPost, accessTokenUrl, nil)
	accessTokenRequest.Header.Set(headerAcceptParam, headerAcceptValue)
	authTokenResponse, err := httpClient.Do(accessTokenRequest)

	if err != nil {
		return OAuthAccessResponse{}, fmt.Errorf("::: could not send HTTP request: %v", err)
	}

	return decodeAccessTokenResponse(authTokenResponse)
}

func decodeAccessTokenResponse(res *http.Response) (OAuthAccessResponse, error) {
	var authResponse OAuthAccessResponse

	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return authResponse, fmt.Errorf("::: could not parse JSON response: %v", err)
	}

	defer res.Body.Close()
	return authResponse, nil
}

func getAccessTokenUrl(code string) string {
	authCallback := viper.GetString(authCallbackUrl)
	authUrl := viper.GetString(authAccessTokenUrl)
	authUrl = fmt.Sprintf(authUrl, credentials.clientId, credentials.clientSecret, code, authCallback)
	return authUrl
}

func getCredentials() *ghCredentials {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("::: Error loading configuration (.env) file")
	}

	return &ghCredentials{
		clientId:     viper.GetString(clientId),
		clientSecret: viper.GetString(clientSecret),
	}
}
