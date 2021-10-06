package main

import (
	"github.com/htenjo/gh_statistics/config"
	"github.com/htenjo/gh_statistics/web"
	_ "github.com/lib/pq"
)

func main() {
	config.InitConfig(".")
	store := config.InitRepository()
	defer store.Close()
	web.Init(store)
}