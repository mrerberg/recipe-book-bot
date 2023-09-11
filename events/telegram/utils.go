package telegram

import (
	tgClient "recipe-book-bot/clients/telegram"
)

func NewInlineKeyboardRow(buttons ...tgClient.InlineKeyboardButton) []tgClient.InlineKeyboardButton {
	var row []tgClient.InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

func NewInlineKeyboardButton(text string, data string) tgClient.InlineKeyboardButton {
	return tgClient.InlineKeyboardButton{Text: text, CallbackData: data}
}
