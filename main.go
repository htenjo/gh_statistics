package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/storage"
	"github.com/htenjo/gh_statistics/web/handler"
	"github.com/htenjo/gh_statistics/web/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	store := initStorage()
	defer store.Close()
	initGin(store)
}

func initGin(store *storage.Storage) {
	authHandler := middleware.NewAuthHandler(store)
	repoHandler := handler.NewRepoHandler(store)

	router := gin.Default()
	router.LoadHTMLGlob("web/template/*")

	router.GET("/ping", middleware.Ping)
	router.GET("/callback", authHandler.CallbackHandler)

	authorized := router.Group("/", middleware.Authenticate())

	{
		authorized.GET("/", handler.IndexHandler)
		authorized.GET(handler.ReposPath, repoHandler.ListRepos)
		authorized.POST(handler.ReposPath, repoHandler.CreateRepos)
	}

	router.Run(":8080")
}

func initStorage() *storage.Storage {
	db, err := sql.Open("sqlite3", "./gh.db")

	if err != nil {
		panic(err)
	}

	return storage.NewStorage(db)
}
