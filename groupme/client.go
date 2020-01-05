package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	c           *http.Client
	BotId       string
	AccessToken string

	BaseUrl string
}

type botPost struct {
	BotId string `json:"bot_id"`
	Text  string `json:"text"`
}

func (c *Client) PostBotMessage(text string) error {
	p := &botPost{
		BotId: c.BotId,
		Text:  text,
	}

	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseUrl+"/bots/post", bytes.NewBuffer(js))
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

func (c *Client) GetGroupMessages(groupId string, params GroupMessageParams) ([]*Message, error) {
	u := fmt.Sprintf("%s/groups/%s/messages", c.BaseUrl, groupId)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("token", c.AccessToken)
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
		q.Set("limit", string(params.Limit))
	} else {
		q.Set("limit", string(100))
	}

	req.URL.RawQuery = q.Encode()
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ms []*Message
	err = json.NewDecoder(resp.Body).Decode(ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}
