package event_consumer

import (
	"context"
	"log"
	"recipe-book-bot/events"
	"time"
)

type Consumer struct {
	batchSize int
	fetcher   events.Fetcher
	processor events.Processor
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		events, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(events) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(context.Background(), events); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
