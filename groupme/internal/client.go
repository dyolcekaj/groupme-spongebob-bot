package internal

import (
	"encoding/json"
	"net/http"
)

// Client for meta calls the bot might need
type Client interface {
	ListBots() ([]Bot, error)
}

type client struct {
	c           *http.Client
	accessToken string

	baseURL string
}

// NewClient for making calls to GroupMe API
func NewClient(baseURL string, accessToken string) Client {
	return &client{
		c:           &http.Client{},
		accessToken: accessToken,
		baseURL:     baseURL,
	}
}

// Bot is a GroupMe bot defined by a user
type Bot struct {
	BotID   string `json:"bot_id"`
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
}

func (c *client) ListBots() ([]Bot, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/bots", nil)
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
