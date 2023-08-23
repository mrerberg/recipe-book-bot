package telegram

import (
	"context"
	"errors"
	"recipe-book-bot/clients/telegram"
	"recipe-book-bot/events"
	lib "recipe-book-bot/lib/error"
	"recipe-book-bot/storage"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Metadata struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.FetchUpdates(p.offset, limit)
	if err != nil {
		return nil, lib.WrapErr("can't fetch events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return lib.WrapErr("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return lib.WrapErr("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return lib.WrapErr("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Metadata, error) {
	meta, ok := event.Metadata.(Metadata)
	if !ok {
		return Metadata{}, lib.WrapErr("can't get meta", ErrUnknownMetaType)
	}

	return meta, nil
}

func event(update telegram.Update) events.Event {
	if update.Message == nil {
		return events.Event{
			Type: events.Unknown,
			Text: "",
		}
	}

	return events.Event{
		Type: events.Message,
		Text: update.Message.Text,
		Metadata: Metadata{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		},
	}
}
