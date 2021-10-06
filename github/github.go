package github

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/config"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	headerAcceptParam      = "Accept"
	headerContentTypeParam = "Content-Type"
	headerAuthorization    = "Authorization"
	headerJsonValue        = "application/json"
	authCodeParam          = "code"
)

var httpClient = http.Client{}

func AuthorizationUrl() string {
	//TODO: includes the state param to avoid CSRF
	log.Println("::: Redirect to GitHub to authorize")
	return fmt.Sprintf(config.GhAuthorizeUrl(), config.GhClientId())
}

func GetUserInfo(authCredentials OAuthCredentials) (GhUser, error) {
	userResponse, err := authGetRequest(authCredentials, config.GhUserApiUrl())

	if err != nil {
		return GhUser{}, fmt.Errorf("::: Error in HTTP request: %v", err)
	}

	var user GhUser
	decodeJsonResponse(userResponse, &user)
	return user, nil
}

func GetOpenPRs(repoName, accessToken string, channel chan RepoPRResponse) {
	var response RepoPRResponse

	if strings.TrimSpace(repoName) == "" {
		response.Error = fmt.Errorf("invalid repo name")
		channel <- response
		return
	}

	repoUrl := config.GhApiBase() + strings.TrimSpace(repoName) + "/pulls?state=open&sort=updated"
	var openPullRequests []PullRequestDetail
	jsonRequest(repoUrl, accessToken, &openPullRequests)
	assignPrOpenFlags(&openPullRequests)
	response.Repo = RepoPR{
		RepositoryName: repoName,
		RepositoryURL:  repoUrl,
		Prs:            openPullRequests,
	}

	channel <- response
}

func Authorize(c *gin.Context) (OAuthCredentials, error) {
	code := c.Query(authCodeParam)
	accessTokenUrl := getAccessTokenUrl(code)
	accessTokenRequest, _ := http.NewRequest(http.MethodPost, accessTokenUrl, nil)
	accessTokenRequest.Header.Set(headerAcceptParam, headerJsonValue)
	authTokenResponse, err := httpClient.Do(accessTokenRequest)

	if err != nil {
		return OAuthCredentials{}, fmt.Errorf("::: Could not send HTTP request: %v", err)
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
	authCallback := config.GhCallbackUrl()
	authUrl := config.GhAccessTokenUrl()
	authUrl = fmt.Sprintf(authUrl, config.GhClientId(), config.GhClientSecret(), code, authCallback)
	log.Printf("::: accessTokenUrl = [%s]", authUrl)
	return authUrl
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
