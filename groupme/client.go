package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Client for making calls to GroupMe API by a
// Command developer
type Client interface {
	PostBotMessage(text string) error
	GetGroupMessages(groupID string, params GroupMessageParams) ([]Message, error)
}

type client struct {
	c           *http.Client
	botID       string
	accessToken string

	baseURL string
}

// NewClient returns a configured Client
func NewClient(botID string, baseURL string, accessToken string) Client {
	return &client{
		c:           &http.Client{},
		botID:       botID,
		accessToken: accessToken,
		baseURL:     baseURL,
	}
}

func (c *client) BotID() string {
	return c.botID
}

func (c *client) AccessToken() string {
	return c.accessToken
}

type botPost struct {
	BotID string `json:"bot_id"`
	Text  string `json:"text"`
}

func (c *client) PostBotMessage(text string) error {
	p := &botPost{
		BotID: c.botID,
		Text:  text,
	}

	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/bots/post", bytes.NewBuffer(js))
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

// GroupMessageParams for searching messages in a chat
type GroupMessageParams struct {
	BeforeID string
	SinceID  string
	AfterID  string
	Limit    int
}

type groupMessageSearchResult struct {
	Response struct {
		Count    int       `json:"count"`
		Messages []Message `json:"messages"`
	} `json:"response"`
}

func (c *client) GetGroupMessages(groupID string, params GroupMessageParams) ([]Message, error) {
	u := fmt.Sprintf("%s/groups/%s/messages", c.baseURL, groupID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("token", c.accessToken)
	if len(params.AfterID) > 0 {
		q.Set("after_id", params.AfterID)
	}
	if len(params.BeforeID) > 0 {
		q.Set("before_id", params.BeforeID)
	}
	if len(params.SinceID) > 0 {
		q.Set("since_id", params.SinceID)
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
