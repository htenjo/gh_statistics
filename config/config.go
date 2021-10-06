package config

import (
	"database/sql"
	"github.com/htenjo/gh_statistics/repository"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

const (
	dbDriver = "postgres"
	dbUrl    = "DATABASE_URL"

	ghClientId       = "GH_CLIENT_ID"
	ghClientSecret   = "GH_CLIENT_SECRET"
	ghAuthorizeUrl   = "GH_AUTHORIZE_URL"
	ghUserApiUrl     = "GH_API_USER_URL"
	ghApiBase        = "GH_API_REPO_URL"
	ghCallbackUrl    = "GH_AUTH_CALLBACK_URL"
	ghAccessTokenUrl = "GH_ACCESS_TOKEN_URL"

	slackWebhookUsePrivate = "SLACK_WEBHOOK_USE_PRIVATE"
	slackPrivateWebhookUrl = "SLACK_PRIVATE_WEBHOOK_URL"
	slackBackendWebhookUlr = "SLACK_BACKEND_WEBHOOK_URL"

	webPort = "PORT"
)

func InitConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("::: Config file not found - Taking values from ENV Variables...")
		} else {
			log.Fatalf("::: Error trying to configure the application - %v", err)
		}
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
	return getValue(dbUrl)
}

func GhClientId() string {
	return getValue(ghClientId)
}

func GhClientSecret() string {
	return getValue(ghClientSecret)
}

func GhAuthorizeUrl() string {
	return getValue(ghAuthorizeUrl)
}

func GhUserApiUrl() string {
	return getValue(ghUserApiUrl)
}

func GhApiBase() string {
	return getValue(ghApiBase)
}

func GhCallbackUrl() string {
	return getValue(ghCallbackUrl)
}

func GhAccessTokenUrl() string {
	return getValue(ghAccessTokenUrl)
}

func SlackWebhookUrl() string {
	if getValueBool(slackWebhookUsePrivate) {
		return getValue(slackPrivateWebhookUrl)
	}

	return getValue(slackBackendWebhookUlr)
}

func WebPort() int {
	return getValueInt(webPort)
}

func getValue(key string) string {
	//return viper.GetString(ghAuthorizeUrl)
	return os.Getenv(key)
}

func getValueInt(key string) int {
	//return viper.GetInt("PORT")
	value := os.Getenv(key)
	intValue, err := strconv.Atoi(value)

	if err != nil {
		log.Printf("::: Error getting config int value [%s] - %v", key, err)
		return -1
	}

	return intValue
}

func getValueBool(key string) bool {
	stringValue := getValue(key)
	result, err := strconv.ParseBool(stringValue)

	if err != nil {
		log.Printf("::: Error getting config bool value [%s] - %v", key, err)
		return false
	}

	return result
}
