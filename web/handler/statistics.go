package handler

import (
	"fmt"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/repository"
	"net/http"
)

type StatsHandler struct {
	storage *repository.UserRepository
}

func NewStatsHandler(storage *repository.UserRepository) StatsHandler {
	return StatsHandler{
		storage: storage,
	}
}

func (s *StatsHandler) Handler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(definition.CookieSession)

	if err != nil {
		fmt.Fprintf(w, "No session found... <a href='/'>Go Home</a>")
	}

	sessionId := cookie.Value
	user, err := s.storage.Find(sessionId)
	fmt.Fprintf(w, "Welcome %s <%s>", user.Username, user.Email)
}
