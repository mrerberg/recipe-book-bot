package telegram

import (
	"encoding/json"
	"fmt"
	"io"
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

func (c *Client) InitBotCommands() error {
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
				Command:     "all",
				Description: "Get all recipes in book",
			},
			{
				Command:     "template",
				Description: "Get template to add recipe",
			},
		},
	}

	err := c.SetCommands(commands)
	if err != nil {
		return err
	}

	return nil
}

// Implements https://core.telegram.org/bots/api#getupdates
func (c *Client) FetchUpdates(offset int, limit int) ([]Update, error) {
	q := url.Values{}

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.makeRequest(getUpdates, q)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// Implements https://core.telegram.org/bots/api#sendmessage
func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}

	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.makeRequest(sendMessage, q)
	if err != nil {
		return lib.WrapErr("could not send message", err)
	}

	return nil
}

// Implements https://core.telegram.org/bots/api#setmycommands
func (c *Client) SetCommands(commands Commands) error {
	q := url.Values{}

	cmds, err := json.Marshal(commands.Commands)
	q.Add("commands", string(cmds))
	q.Add("language_code", commands.LanguageCode)

	_, err = c.makeRequest(setMyCommands, q)
	if err != nil {
		return lib.WrapErr("could not send message", err)
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) makeRequest(method string, query url.Values) ([]byte, error) {
	const errMsg = "could not make request"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// TODO: read errors.Is, errors.As
		// TODO: Handle error! Если приходит 400, приложение не сигнализирует об этом
		fmt.Printf("[ERROR] [makeRequest] %v \n", err)
		return nil, lib.WrapErr(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, lib.WrapErr(errMsg, err)
	}

	fmt.Printf("SAS %v", req.URL)

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, lib.WrapErr(errMsg, err)
	}

	return body, nil
}
