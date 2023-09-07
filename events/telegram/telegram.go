package telegram

import (
	"context"
	"errors"
	"recipe-book-bot/clients/telegram"
	"recipe-book-bot/events"
	lib "recipe-book-bot/lib/error"
	"recipe-book-bot/storage"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
	cache   *memcache.Client
}

type Metadata struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage, cache *memcache.Client) *Processor {
	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
		cache:   cache,
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
	case events.CallBack:
		return p.processCallBack(ctx, event)
	case events.Unknown:
		return lib.WrapErr("can't process message", ErrUnknownEventType)
	default:
		return lib.WrapErr("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	msg := "can't process message"
	if err != nil {
		return lib.WrapErr(msg, err)
	}

	if err = p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return lib.WrapErr(msg, err)
	}

	return nil
}

func (p *Processor) processCallBack(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	msg := "can't process callback"
	if err != nil {
		return lib.WrapErr(msg, err)
	}

	if err = p.doCallBack(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return lib.WrapErr(msg, err)
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
	if update.Message == nil && update.CallbackQuery == nil {
		return events.Event{
			Type: events.Unknown,
			Text: "",
		}
	}

	if update.CallbackQuery != nil {
		return events.Event{
			Type: events.CallBack,
			Text: update.CallbackQuery.Data,
			Metadata: Metadata{
				ChatID:   update.CallbackQuery.Message.Chat.ID,
				Username: update.CallbackQuery.Message.Chat.Username,
			},
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
