package config

import (
	"log"
	"os"
)

type Config struct {
	TgBotToken            string
	TgAPIHost             string
	MongoConnectionString string
	CacheHost             string
}

func MustLoad() Config {
	token := os.Getenv("TG_BOT_TOKEN")
	mongoConnectionString := os.Getenv("DB_STRING")
	cacheHost := os.Getenv("CACHE_HOST")
	apiHost := "api.telegram.org"

	if token == "" {
		log.Fatal("TG_BOT_TOKEN variable was not specified")
	}

	if apiHost == "" {
		log.Fatal("tg-api-host variable was not specified")
	}

	if mongoConnectionString == "" {
		log.Fatal("DB_STRING variable was not specified")
	}

	if cacheHost == "" {
		log.Fatal("CACHE_HOST variable was not specified")
	}

	return Config{
		TgBotToken:            token,
		TgAPIHost:             apiHost,
		MongoConnectionString: mongoConnectionString,
		CacheHost:             cacheHost,
	}
}
