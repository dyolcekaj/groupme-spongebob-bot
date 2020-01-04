package groupme

import "fmt"

type Post struct {
	BotId string `json:"bot_id"`
	Test  string `json:"text"`
}

/*
{
  "attachments": [],
  "avatar_url": "https://i.groupme.com/123456789",
  "created_at": 1302623328,
  "group_id": "1234567890",
  "id": "1234567890",
  "name": "John",
  "sender_id": "12345",
  "sender_type": "user",
  "source_guid": "GUID",
  "system": false,
  "text": "Hello world ☃☃",
  "user_id": "1234567890"
}
*/

type SenderType string

const (
	UserSender SenderType = "user"
	BotSender  SenderType = "bot"
)

const (
	url string = "https://api.groupme.com/v3/groups/%d/messages?token=%s"
)

type Message struct {
	CreatedAt  int        `json:"created_at"`
	GroupId    int        `json:"group_id"`
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	SenderId   int        `json:"sender_id"`
	SenderType SenderType `json:"sender_type"`
	SourceGuid string     `json:"source_guid"`
	System     bool       `json:"system"`
	Text       string     `json:"text"`
	UserId     int        `json:"user_id"`
}

func (m *Message) URL(token string) string {
	return fmt.Sprintf(url, m.GroupId, token)
}
