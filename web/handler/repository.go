package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/config"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/github"
	"github.com/htenjo/gh_statistics/repository"
	"github.com/htenjo/gh_statistics/slack"
	"log"
	"net/http"
	"strings"
)

const ReposPath = "/repos"
const ReposOpenPRNotification = "/repos/open-pr/notification"

type RepoHandler struct {
	store *repository.UserRepository
}

func NewRepoHandler(store *repository.UserRepository) RepoHandler {
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

	_, err := h.store.UpdateGitRepositories(sessionId, strings.Join(repoUrls, ","))

	if err != nil {
		log.Printf("::: Error updating git repositories, %v", err)
	}

	c.Redirect(http.StatusFound, ReposPath)
}

func (h *RepoHandler) SendOpenPRNotification(c *gin.Context) {
	info := h.getOpenPRInformation(c)
	slack.SendSlackMessage("Open PRs", &info)
}

// SendPRNotification Temp handler for easy job notifications
func (h *RepoHandler) SendPRNotification(c *gin.Context) {
	sessionId := c.GetHeader("x-session-id")
	authToken := c.GetHeader("x-cron-token")

	if sessionId == "" {
		c.JSON(400, gin.H{"message": "Required header not found"})
		return
	}

	if authToken != config.AuthTempToken() {
		c.JSON(401, gin.H{"message": "User not authorized"})
		return
	}

	c.Set(definition.SessionId, sessionId)
	h.SendOpenPRNotification(c)
}

func (h *RepoHandler) getOpenPRInformation(c *gin.Context) []github.RepoPR {
	sessionId := c.GetString(definition.SessionId)
	user, _ := h.store.Find(sessionId)
	repos := strings.Split(user.Repos, ",")

	info := make([]github.RepoPR, 0)
	prChannel := make(chan github.RepoPRResponse)

	for _, repoName := range repos {
		go github.GetOpenPRs(repoName, user.AccessToken, prChannel)
	}

	for i := 0; i < len(repos); i++ {
		repoResponse := <-prChannel

		if repoResponse.Error != nil {
			log.Print(repoResponse.Error)
			continue
		}

		info = append(info, repoResponse.Repo)
	}

	return info
}
