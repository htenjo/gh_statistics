package github

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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
	userApiUrl         = "GH_USER_URL"
)

var credentials = getAppCredentials()
var httpClient = http.Client{}

func AuthorizationUrl() string {
	//TODO: includes the state param to avoid CSRF
	return fmt.Sprintf(viper.GetString(authorizeUrl), credentials.clientId)
}

func Authorize(c *gin.Context) (OAuthCredentials, error) {
	code := c.Query(authCodeParam)
	accessTokenUrl := getAccessTokenUrl(code)
	accessTokenRequest, _ := http.NewRequest(http.MethodPost, accessTokenUrl, nil)
	accessTokenRequest.Header.Set(headerAcceptParam, headerAcceptValue)
	authTokenResponse, err := httpClient.Do(accessTokenRequest)

	if err != nil {
		return OAuthCredentials{}, fmt.Errorf("::: could not send HTTP request: %v", err)
	}

	var authResponse OAuthCredentials
	decodeJsonResponse(authTokenResponse, &authResponse)
	return authResponse, nil
}

func decodeJsonResponse(res *http.Response, fillIn interface{}) error {
	if err := json.NewDecoder(res.Body).Decode(fillIn); err != nil {
		return fmt.Errorf("::: could not parse JSON response: %v", err)
	}

	defer res.Body.Close()
	return nil
}

func getAccessTokenUrl(code string) string {
	authCallback := viper.GetString(authCallbackUrl)
	authUrl := viper.GetString(authAccessTokenUrl)
	authUrl = fmt.Sprintf(authUrl, credentials.clientId, credentials.clientSecret, code, authCallback)
	return authUrl
}

func getAppCredentials() *ghCredentials {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("::: Error loading configuration (.env) file")
	}

	return &ghCredentials{
		clientId:     viper.GetString(clientId),
		clientSecret: viper.GetString(clientSecret),
	}
}

func GetUserInfo(authCredentials OAuthCredentials) (GhUser, error) {
	userResponse, err := authGetRequest(authCredentials, viper.GetString(userApiUrl))

	if err != nil {
		return GhUser{}, fmt.Errorf("::: Error in HTTP request: %v", err)
	}

	var user GhUser
	decodeJsonResponse(userResponse, &user)
	return user, nil
}

func authGetRequest(authCredentials OAuthCredentials, url string) (*http.Response, error) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", authCredentials.getFullTokenHeader())
	response, err := httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("::: Error in HTTP request: %v", err)
	}

	return response, nil
}
