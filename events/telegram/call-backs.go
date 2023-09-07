package telegram

import (
	"context"
	"log"
	"recipe-book-bot/events"
	"strings"
)

const (
	DeleteCb = "delete"
	GetCb    = "get"
)

func (p *Processor) doCallBack(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	handleUnknownCb := func() error {
		log.Printf("got new callback '%s' from '%s", text, username)
		return p.tg.SendMessage(chatID, unknownCbMsg)
	}

	if !events.IsCallBack(text) {
		return handleUnknownCb()
	}

	parts := events.ParseCallBack(text)
	cbType := parts[1]
	recipeName := parts[len(parts)-1]

	switch cbType {
	case DeleteCb:
		return p.deleteRecipe(ctx, chatID, recipeName, username)
	case GetCb:
		return p.getFullRecipe(ctx, chatID, recipeName, username)
	default:
		return handleUnknownCb()
	}
}
