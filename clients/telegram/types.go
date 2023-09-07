package telegram

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Commands struct {
	LanguageCode string    `json:"language_code"`
	Commands     []Command `json:"commands"`
}

type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type CallbackQuery struct {
	ID           string  `json:"id"`
	From         From    `json:"from"`
	Message      Message `json:"message"`
	ChatInstance string  `json:"chat_instance"`
	Data         string  `json:"data"`
}

type Message struct {
	MessageID   int    `json:"message_id"`
	From        From   `json:"from"`
	Chat        Chat   `json:"chat"`
	Date        int    `json:"date"`
	Text        string `json:"text"`
	ReplyMarkup struct {
		InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	} `json:"reply_markup"`
}
