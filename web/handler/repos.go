package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/storage"
	"net/http"
)

const ReposPath = "/repos"

type RepoHandler struct {
	store *storage.Storage
}

func NewRepoHandler(store *storage.Storage) RepoHandler {
	return RepoHandler{store: store}
}

func (h *RepoHandler) ListRepos(c *gin.Context) {
	sessionId := c.GetString(definition.SessionId)
	user, _ := h.store.Find(sessionId)

	c.HTML(http.StatusOK, "repos.html", gin.H{
		"title": "Repositories",
		"repos": user.Repos,
	})
}

func (h *RepoHandler) CreateRepos(c *gin.Context) {
	sessionId := c.GetString(definition.SessionId)
	repoUrls := c.Request.FormValue("repoUrls")
	h.store.UpdateRepos(sessionId, repoUrls)
	c.Redirect(http.StatusFound, ReposPath)
}
