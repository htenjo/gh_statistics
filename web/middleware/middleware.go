package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/htenjo/gh_statistics/definition"
	"github.com/htenjo/gh_statistics/github"
	"log"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(definition.CookieSession)

		if err != nil || cookie == "" {
			log.Println("::: Authentication information NOT FOUND")
			c.Redirect(http.StatusTemporaryRedirect, github.AuthorizationUrl())
			c.AbortWithStatus(http.StatusTemporaryRedirect)
		} else {
			c.Set(definition.SessionId, cookie)
		}
	}
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
