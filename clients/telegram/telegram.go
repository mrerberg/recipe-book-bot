package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	lib "recipe-book-bot/lib/error"
	"strconv"
)

const (
	getUpdates    = "getUpdates"
	sendMessage   = "sendMessage"
	setMyCommands = "setMyCommands"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(token string, host string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) InitBotCommands() {
	var commands = Commands{
		LanguageCode: "en",
		Commands: []Command{
			{
				Command:     "start",
				Description: "start",
			},
			{
				Command:     "help",
				Description: "Display help information",
			},
			{
				Command:     "add",
				Description: "Start adding new recipe",
			},
			{
				Command:     "all",
				Description: "Get all recipes in book",
			},
		},
	}

	err := c.SetCommands(commands)
	if err != nil {
		log.Printf("could not init bot commands. Error: %v", err)
	}
}

// Implements https://core.telegram.org/bots/api#getupdates
func (c *Client) FetchUpdates(offset int, limit int) ([]Update, error) {
	q := url.Values{}

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	var res UpdateResponse
	err := c.makeRequest(getUpdates, q, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

// Implements https://core.telegram.org/bots/api#sendmessage
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}

	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	err := c.makeRequest(sendMessage, q, nil)
	if err != nil {
		return lib.WrapErr("could not send message", err)
	}

	return nil
}

func (c *Client) SendMessageWithMarkup(chatID int, text string, inlineKeyboard *InlineKeyboard) error {
	q := url.Values{}

	data, _ := json.Marshal(inlineKeyboard)

	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("reply_markup", string(data))

	err := c.makeRequest(sendMessage, q, nil)
	if err != nil {
		return lib.WrapErr("could not send message", err)
	}

	return nil
}

// Implements https://core.telegram.org/bots/api#setmycommands
func (c *Client) SetCommands(commands Commands) error {
	q := url.Values{}

	cmds, err := json.Marshal(commands.Commands)
	if err != nil {
		return err
	}
	q.Add("commands", string(cmds))
	q.Add("language_code", commands.LanguageCode)

	err = c.makeRequest(setMyCommands, q, nil)
	if err != nil {
		return lib.WrapErr("could not send message", err)
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) makeRequest(method string, query url.Values, v interface{}) error {
	const errMsg = "could not make request"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// TODO: read errors.Is, errors.As
		return lib.WrapErr(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return lib.WrapErr(errMsg, err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return lib.WrapErr(errMsg, err)
	}

	if res.StatusCode != http.StatusOK {
		return lib.WrapErr(fmt.Sprintf("Request status: %v", res.StatusCode), nil)
	}

	if v == nil {
		return nil
	}

	if err = json.Unmarshal(body, v); err != nil {
		return lib.WrapErr("Error: could not unmarshal JSON response", err)
	}

	return nil
}
