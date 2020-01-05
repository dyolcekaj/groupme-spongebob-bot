package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	objects "github.com/dyolcekaj/groupme-spongebob-bot/groupme/internal"
	"net/http"
)

type Client struct {
	c           *http.Client
	BotId       string
	AccessToken string

	BaseUrl string
}

func (c *Client) PostBotMessage(text string) error {
	p := &objects.Post{
		BotId: c.BotId,
		Text:  text,
	}

	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseUrl + "/bots/post", bytes.NewBuffer(js))
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
