package config

import (
	"flag"
	"log"
)

type Config struct {
	TgBotToken            string
	TgApiHost             string
	MongoConnectionString string
}

func MustLoad() Config {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	apiHost := flag.String("tg-api-host", "api.telegram.org", "host of bot api")
	mongoConnectionString := flag.String("mongo-connection-string", "", "connection string for MongoDB")

	flag.Parse()

	if *token == "" {
		log.Fatal("tg-bot-token flag was not specified")
	}

	if *apiHost == "" {
		log.Fatal("tg-api-host flag was not specified")
	}

	if *mongoConnectionString == "" {
		log.Fatal("mongo-connection-string flag was not specified")
	}

	return Config{
		TgBotToken:            *token,
		TgApiHost:             *apiHost,
		MongoConnectionString: *mongoConnectionString,
	}
}
