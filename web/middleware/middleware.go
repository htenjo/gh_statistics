package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/github"
	"github.com/htenjo/gh_statistics/repository"
	"log"
	"net/http"
)

func Authenticate(store *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(definition.CookieSession)

		if err != nil || cookie == "" {
			log.Println("::: Authentication information NOT FOUND")
			redirectGh(c)
		} else {
			_, err := store.Find(cookie)

			if err != nil {
				log.Println("::: Session found but user not found !!!")
				redirectGh(c)
			}

			//log.Println("::: Authentication OK")
			c.Set(definition.SessionId, cookie)
		}
	}
}

func redirectGh(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, github.AuthorizationUrl())
	c.AbortWithStatus(http.StatusTemporaryRedirect)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}