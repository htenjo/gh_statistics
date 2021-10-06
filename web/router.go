package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/repository"
	"github.com/htenjo/gh_statistics/web/handler"
	"github.com/htenjo/gh_statistics/web/middleware"
	"github.com/spf13/viper"
)

func Init(store *repository.UserRepository) {
	authHandler := middleware.NewAuthHandler(store)
	repoHandler := handler.NewRepoHandler(store)

	router := gin.Default()
	router.LoadHTMLGlob("web/template/*")

	router.GET("/ping", middleware.Ping)
	router.GET("/callback", authHandler.CallbackHandler)
	authGuard := middleware.Authenticate(store)

	router.GET("/", authGuard, handler.IndexHandler)
	router.GET(handler.ReposPath, authGuard, repoHandler.ListRepos)
	router.POST(handler.ReposPath, authGuard, repoHandler.CreateRepos)
	router.POST(handler.ReposOpenPRNotification, authGuard, repoHandler.SendOpenPRNotification)

	if err := router.Run(port()); err != nil {
		panic(err)
	}
}

func port() string {
	return fmt.Sprintf(":%d", viper.GetInt("PORT"))
}