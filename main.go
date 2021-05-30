package main

import (
	"database/sql"
	"github.com/htenjo/gh_statistics/storage"
	"github.com/htenjo/gh_statistics/web"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	store := initStorage()
	defer store.Close()
	authHandler := web.NewAuthHandler(store)
	statsHandler := web.NewStatsHandler(store)

	http.HandleFunc("/", authHandler.LoginHandler)
	http.HandleFunc("/callback", authHandler.CallbackHandler)
	http.HandleFunc("/home", statsHandler.Handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initStorage() *storage.Storage {
	db, err := sql.Open("sqlite3", "./gh.db")
	checkErr(err)
	return storage.NewStorage(db)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
