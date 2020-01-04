package groupme

import (
	"unicode"
)

type Post struct {
	BotId string `json:"bot_id"`
	Text  string `json:"text"`
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

type BotCommand interface {
	CreatePost(botId string) *Post
}

type lastMessageCommand struct {
	m *Message

}

func (c *lastMessageCommand) CreatePost(botId string) *Post {
	return nil // TODO
}

type thisMessageCommand struct {
	m *Message
	parsedText string
}
func (c *thisMessageCommand) CreatePost(botId string) *Post {
	return &Post{botId, translateText(c.parsedText)}
}

func (m *Message) ParseCommand() (BotCommand, error) {
	return &thisMessageCommand{}, nil
}

func translateText(text string) string {
	ret := []rune{}
	i := false
	for _, r := range text {
		if i {
			ret = append(ret, unicode.ToUpper(r))
		} else {
			ret = append(ret, unicode.ToLower(r))
		}

		if r != ' ' {
			i = !i
		}
	}

	return string(ret)
}

