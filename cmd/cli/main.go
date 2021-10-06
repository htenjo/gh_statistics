package main

import (
	"flag"
	"github.com/htenjo/gh_statistics/cli"
	"github.com/htenjo/gh_statistics/config"
)

func main() {
	sid := flag.String("sid", "-1", "This is the SESSION_ID previously configured from web")
	flag.Parse()

	if *sid == "-1" {
		panic("::: Missing SID flag")
	}

	store := config.InitRepository()
	repoHandler := cli.NewHandler(store)
	repoHandler.SendOpenPRNotification(*sid)
}