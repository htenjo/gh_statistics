package main

import (
	"github.com/htenjo/gh_statistics/web"
	"log"
	"net/http"
)

func main() {
	initRouter()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initRouter() {
	http.HandleFunc("/", web.LoginHandler)
	http.HandleFunc("/callback", web.CallbackHandler)
	http.HandleFunc("/home", web.Home)
}
