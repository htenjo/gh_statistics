package web

import (
	"fmt"
	"github.com/htenjo/gh_statistics/github"
	"github.com/htenjo/gh_statistics/storage"
	"net/http"
	"strconv"
)

const CookieSession = "gh_session"

type AuthHandler struct {
	storage *storage.Storage
}

func NewAuthHandler(storage *storage.Storage) AuthHandler {
	return AuthHandler{storage: storage}
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: includes the state param to avoid CSRF
	if a.sessionCookieExist(r) {
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, github.AuthorizationUrl(), http.StatusTemporaryRedirect)
	}
}

func (a *AuthHandler) CallbackHandler(response http.ResponseWriter, request *http.Request) {
	credentials, err := github.Authorize(request)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, err.Error())
	}

	ghUser, err := github.GetUserInfo(credentials)
	user, err := a.storage.Find(strconv.FormatInt(ghUser.Id, 10))

	if err != nil {
		user, err = a.persistGhUser(ghUser, credentials)
	}

	persistCookie(response, user.SessionId)
	response.Header().Set("Location", "/home")
	response.WriteHeader(http.StatusFound)
}

func (a *AuthHandler) persistGhUser(ghUser github.GhUser, credentials github.OAuthCredentials) (storage.User, error) {
	user := storage.User{
		SessionId:   strconv.FormatInt(ghUser.Id, 10),
		AccessToken: credentials.AccessToken,
		Email:       ghUser.Email,
		Username:    ghUser.Login,
	}
	return a.storage.Save(user)
}

func persistCookie(response http.ResponseWriter, sessionId string) {
	cookie := http.Cookie{Name: CookieSession, Value: sessionId}
	http.SetCookie(response, &cookie)
}

func (a *AuthHandler) sessionCookieExist(r *http.Request) bool {
	cookie, err := r.Cookie(CookieSession)

	if err != nil || cookie.Value == "" {
		return false
	}

	return true
}
