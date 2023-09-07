package main

import (
	"log"
	tgClient "recipe-book-bot/clients/telegram"
	"recipe-book-bot/config"
	eventConsumer "recipe-book-bot/consumer/event-consumer"
	"recipe-book-bot/events/telegram"
	"recipe-book-bot/storage/mongo"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
)

const (
	batchSize = 100
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[APP] Could not load env file")
	}

	cfg := config.MustLoad()

	client := tgClient.New(cfg.TgBotToken, cfg.TgAPIHost)
	client.InitBotCommands()

	storage := mongo.New(cfg.MongoConnectionString)

	cache := memcache.New(cfg.CacheHost)

	eventsProcessor := telegram.New(
		client,
		storage,
		cache,
	)

	log.Print("[APP] Service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("[APP] Service stopped", err)
	}
}
