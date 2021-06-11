package main

import (
	"database/sql"
	"flag"
	"github.com/htenjo/gh_statistics/cli"
	"github.com/htenjo/gh_statistics/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sid := flag.String("sid", "-1", "This is the SESSION_ID previously configured from web")
	flag.Parse()

	if *sid == "-1" {
		panic("::: Missing SID flag")
	}

	store := initStorage()
	repoHandler := cli.NewHandler(store)
	repoHandler.SendOpenPRNotification(*sid)
}

func initStorage() *storage.Storage {
	db, err := sql.Open("sqlite3", "./gh.db")

	if err != nil {
		panic(err)
	}

	return storage.NewStorage(db)
}
