package internal

import (
	"encoding/json"
	"net/http"
)

type Client interface {
	ListBots() ([]Bot, error)
}

type client struct {
	c           *http.Client
	accessToken string

	baseUrl string
}

func NewClient(baseUrl string, accessToken string) Client {
	return &client{
		c:           &http.Client{},
		accessToken: accessToken,
		baseUrl:     baseUrl,
	}
}

type Bot struct {
	BotId   string `json:"bot_id"`
	GroupId string `json:"group_id"`
	Name    string `json:"name"`
}

func (c *client) ListBots() ([]Bot, error) {
	req, err := http.NewRequest("GET", c.baseUrl+"/bots", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("token", c.accessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = struct {
		Response []Bot
	}{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Response, nil
}
