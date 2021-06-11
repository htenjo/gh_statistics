package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/github"
	"github.com/htenjo/gh_statistics/slack"
	"github.com/htenjo/gh_statistics/storage"
	"net/http"
	"strings"
)

const ReposPath = "/repos"
const ReposOpenPRNotification = "/repos/open-pr/notification"

type RepoHandler struct {
	store *storage.Storage
}

func NewRepoHandler(store *storage.Storage) RepoHandler {
	return RepoHandler{store: store}
}

func (h *RepoHandler) ListRepos(c *gin.Context) {
	info := h.getOpenPRInformation(c)

	c.HTML(http.StatusOK, "repos.html", gin.H{
		"title": "Repositories",
		"info":  info,
	})
}

func (h *RepoHandler) CreateRepos(c *gin.Context) {
	sessionId := c.GetString(definition.SessionId)
	repoUrlsParam := c.Request.FormValue("repoUrls")
	repoUrls := strings.Split(repoUrlsParam, ",")

	for i := 0; i < len(repoUrls); i++ {
		repoUrls[i] = strings.TrimSpace(repoUrls[i])
	}

	h.store.UpdateRepos(sessionId, strings.Join(repoUrls, ","))
	c.Redirect(http.StatusFound, ReposPath)
}

func (h *RepoHandler) SendOpenPRNotification(c *gin.Context) {
	info := h.getOpenPRInformation(c)
	slack.SendSlackMessage("Open PRs", &info)
}

func (h *RepoHandler) getOpenPRInformation(c *gin.Context) []github.RepoPR {
	sessionId := c.GetString(definition.SessionId)
	user, _ := h.store.Find(sessionId)
	repos := strings.Split(user.Repos, ",")

	info := make([]github.RepoPR, 0)
	prChannel := make(chan github.RepoPR)

	for _, repoName := range repos {
		go github.GetOpenPRs(repoName, user.AccessToken, prChannel)
	}

	for i := 0; i < len(repos); i++ {
		repoPR := <-prChannel
		info = append(info, repoPR)
	}

	return info
}
