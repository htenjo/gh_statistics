package web

import (
	"fmt"
	"github.com/htenjo/gh_statistics/github"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: includes the state param to avoid CSRF
	http.Redirect(w, r, github.AuthorizationUrl(), http.StatusTemporaryRedirect)
}

func CallbackHandler(response http.ResponseWriter, request *http.Request) {
	oauthResponse, err := github.Authorize(request)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, err.Error())
	}

	response.Header().Set("Location", "/home?access_token="+oauthResponse.AccessToken)
	response.WriteHeader(http.StatusFound)
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome %s", "-----")
}
