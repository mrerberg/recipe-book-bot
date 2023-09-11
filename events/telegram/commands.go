package telegram

import (
	"context"
	"log"
	"recipe-book-bot/events"
	"strings"
)

const (
	StartCmd   = "/start"
	HelpCmd    = "/help"
	AllCmd     = "/all"
	AllTestCmd = "/all1"
	AddCmd     = "/add"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	if events.IsCommand(text) {
		log.Printf("[COMMANDS] Got new command '%s' from '%s", text, username)

		err := p.cache.Delete(username)
		if err != nil {
			log.Printf("[COMMANDS] Error with cache: %v", err)
		}
	}

	_, err := p.cache.Get(username)

	if !events.IsCommand(text) {
		if err == nil {
			return p.saveRecipe(ctx, chatID, text, username)
		}

		return p.getRecipe(ctx, chatID, text, username)
	}

	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	case AddCmd:
		return p.startRecipeSave(ctx, chatID, username)
	case AllCmd:
		return p.sendAll(ctx, chatID, username, 1, 0) // TODO: fix 0
	default:
		return p.tg.SendMessage(chatID, unknownCommandMsg)
	}
}
