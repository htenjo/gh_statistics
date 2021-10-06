package config

import (
	"database/sql"
	"github.com/htenjo/gh_statistics/repository"
	"github.com/spf13/viper"
	"log"
)

const (
	dbDriver       = "postgres"
	dbUrl          = "DATABASE_URL"
	ghClientId     = "GH_CLIENT_ID"
	ghClientSecret = "GH_CLIENT_SECRET"
)

func InitConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("::: Error loading configuration (config.env) file - %v", err)
	}
}

func InitRepository() *repository.UserRepository {
	db, err := sql.Open(dbDriver, dbConnectionUrl())

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	log.Println("::: Database started...")
	return repository.NewUserRepository(db)
}

func dbConnectionUrl() string {
	return viper.GetString(dbUrl)
}

func GhClientId() string {
	return viper.GetString(ghClientId)
}

func GhClientSecret() string {
	return viper.GetString(ghClientSecret)
}
