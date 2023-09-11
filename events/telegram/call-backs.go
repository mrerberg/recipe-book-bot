package telegram

import (
	"context"
	"log"
	"recipe-book-bot/events"
	"strconv"
	"strings"
)

const (
	DeleteCb = "delete"
	GetCb    = "get"
	GetAllCb = "getall"
)

func (p *Processor) doCallBack(ctx context.Context, text string, chatID int, username string, messageID int) error {
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
	value := parts[len(parts)-1]

	switch cbType {
	case DeleteCb:
		return p.deleteRecipe(ctx, chatID, value, username)
	case GetCb:
		return p.getFullRecipe(ctx, chatID, value, username)
	case GetAllCb:
		page, _ := strconv.ParseInt(value, 10, 64)
		return p.sendAll(ctx, chatID, username, page, messageID)
	default:
		return handleUnknownCb()
	}
}
