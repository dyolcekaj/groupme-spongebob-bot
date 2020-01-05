package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Client interface {
	BotId() string
	AccessToken() string

	PostBotMessage(text string) error
	GetGroupMessages(groupId string, params GroupMessageParams) ([]Message, error)
}

type client struct {
	c           *http.Client
	botId       string
	accessToken string

	baseUrl string
}

func NewClient(botId string, baseUrl string) Client {
	return &client{
		c:           &http.Client{},
		botId:       botId,
		accessToken: "",
		baseUrl:     baseUrl,
	}
}

func NewClientWithToken(botId string, baseUrl string, accessToken string) Client {
	return &client{
		c:           &http.Client{},
		botId:       botId,
		accessToken: accessToken,
		baseUrl:     baseUrl,
	}
}

func (c *client) BotId() string {
	return c.botId
}

func (c *client) AccessToken() string {
	return c.accessToken
}

type botPost struct {
	BotId string `json:"bot_id"`
	Text  string `json:"text"`
}

func (c *client) PostBotMessage(text string) error {
	p := &botPost{
		BotId: c.botId,
		Text:  text,
	}

	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/bots/post", bytes.NewBuffer(js))
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected 201 or 202, got %d - %s", resp.StatusCode, resp.Status)
	}

	return nil
}

type GroupMessageParams struct {
	BeforeId string
	SinceId  string
	AfterId  string
	Limit    int
}

type groupMessageSearchResult struct {
	Response struct {
		Count    int       `json:"count"`
		Messages []Message `json:"messages"`
	} `json:"response"`
}

func (c *client) GetGroupMessages(groupId string, params GroupMessageParams) ([]Message, error) {
	u := fmt.Sprintf("%s/groups/%s/messages", c.baseUrl, groupId)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("token", c.accessToken)
	if len(params.AfterId) > 0 {
		q.Set("after_id", params.AfterId)
	}
	if len(params.BeforeId) > 0 {
		q.Set("before_id", params.BeforeId)
	}
	if len(params.SinceId) > 0 {
		q.Set("since_id", params.SinceId)
	}

	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	} else {
		q.Set("limit", strconv.Itoa(100))
	}

	req.URL.RawQuery = q.Encode()
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result groupMessageSearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Response.Messages, nil
}
