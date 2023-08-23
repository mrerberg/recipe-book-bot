package main

import (
	"log"
	tgClient "recipe-book-bot/clients/telegram"
	"recipe-book-bot/config"
	event_consumer "recipe-book-bot/consumer/event-consumer"
	"recipe-book-bot/events/telegram"
	"recipe-book-bot/storage/mongo"
)

const (
	batchSize = 100
)

func main() {
	cfg := config.MustLoad()

	client := tgClient.New(cfg.TgBotToken, cfg.TgApiHost)
	client.InitBotCommands()

	storage := mongo.New(cfg.MongoConnectionString)

	eventsProcessor := telegram.New(
		client,
		storage,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
