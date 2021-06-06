package github

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	clientId               = "GH_CLIENT_ID"
	clientSecret           = "GH_CLIENT_SECRET"
	authorizeUrl           = "GH_AUTHORIZE_URL"
	authAccessTokenUrl     = "GH_ACCESS_TOKEN_URL"
	authCallbackUrl        = "GH_AUTH_CALLBACK_URL"
	ghHtmlBase             = "GH_HTML_BASE_URL"
	headerAcceptParam      = "Accept"
	headerContentTypeParam = "Content-Type"
	headerAuthorization    = "Authorization"
	headerJsonValue        = "application/json"
	authCodeParam          = "code"
	userApiUrl             = "GH_API_USER_URL"
	ghApiBase              = "GH_API_REPO_URL"
)

var credentials = getAppCredentials()
var httpClient = http.Client{}

func AuthorizationUrl() string {
	//TODO: includes the state param to avoid CSRF
	return fmt.Sprintf(viper.GetString(authorizeUrl), credentials.clientId)
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

func GetOpenPRs(repoName, accessToken string, channel chan RepoPR) {
	repoUrl := viper.GetString(ghApiBase) + strings.TrimSpace(repoName) + "/pulls?state=open&sort=updated"
	var openPullRequests []PullRequestDetail
	jsonRequest(repoUrl, accessToken, &openPullRequests)
	assignPrOpenFlags(&openPullRequests)
	channel <- RepoPR{
		Repository: repoName,
		Prs:        openPullRequests,
	}
}

func Authorize(c *gin.Context) (OAuthCredentials, error) {
	code := c.Query(authCodeParam)
	accessTokenUrl := getAccessTokenUrl(code)
	accessTokenRequest, _ := http.NewRequest(http.MethodPost, accessTokenUrl, nil)
	accessTokenRequest.Header.Set(headerAcceptParam, headerJsonValue)
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

func authGetRequest(authCredentials OAuthCredentials, url string) (*http.Response, error) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(headerAuthorization, authCredentials.getFullTokenHeader())
	response, err := httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("::: Error in HTTP request: %v", err)
	}

	return response, nil
}

func jsonRequest(url, accessToken string, target interface{}) error {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(headerContentTypeParam, headerJsonValue)
	request.Header.Set(headerAuthorization, "Bearer "+accessToken)
	response, err := httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("::: Error in HTTP request: %v", err)
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

func assignPrOpenFlags(pullRequests *[]PullRequestDetail) {
	currentTime := time.Now()

	for i := 0; i < len(*pullRequests); i++ {
		pr := &(*pullRequests)[i]
		openHours := currentTime.Sub(pr.CreatedAt).Hours()

		if openHours < 4 {
			pr.ReviewFlag = Green
		} else if 4 <= openHours && openHours < 8 {
			pr.ReviewFlag = Yellow
		} else {
			pr.ReviewFlag = Red
		}
	}
}
