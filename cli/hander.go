package cli

import (
	"github.com/htenjo/gh_statistics/github"
	"github.com/htenjo/gh_statistics/slack"
	"github.com/htenjo/gh_statistics/storage"
	"strings"
)

type Handler struct {
	store *storage.Storage
}

func NewHandler(store *storage.Storage) Handler {
	return Handler{store: store}
}

func (h *Handler) SendOpenPRNotification(sessionId string) {
	info := h.getOpenPRInformation(sessionId)
	slack.SendSlackMessage("Open PRs", &info)
}

func (h *Handler) getOpenPRInformation(sessionId string) []github.RepoPR {
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
