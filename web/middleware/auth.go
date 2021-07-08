package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/github"
	"github.com/htenjo/gh_statistics/storage"
	"log"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	storage *storage.Storage
}

func NewAuthHandler(storage *storage.Storage) AuthHandler {
	return AuthHandler{storage: storage}
}

func (a *AuthHandler) CallbackHandler(c *gin.Context) {
	log.Println("::: Handling Github Callback")
	credentials, err := github.Authorize(c)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	ghUser, err := github.GetUserInfo(credentials)
	user, err := a.storage.Find(strconv.FormatInt(ghUser.Id, 10))

	if err != nil {
		user, err = a.persistGhUser(ghUser, credentials)
	}

	persistCookie(c, user.SessionId)
	c.Redirect(http.StatusFound, "/")
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

func persistCookie(c *gin.Context, sessionId string) {
	c.SetCookie(definition.CookieSession, sessionId, 0, "/", "localhost", true, true)
}
