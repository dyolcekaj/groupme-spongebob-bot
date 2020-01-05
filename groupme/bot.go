package groupme

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SenderType string
type AttachmentType string

const (
	UserSender SenderType = "user"
	BotSender  SenderType = "bot"

	ImageType    AttachmentType = "image"
	LocationType AttachmentType = "location"
	MentionType  AttachmentType = "mentions"
	EmojiType    AttachmentType = "emoji"
	SplitType    AttachmentType = "split"

	DefaultUrl = "https://api.groupme.com/v3"
)

type Message struct {
	Attachments []Attachment `json:"attachments"`
	CreatedAt   int          `json:"created_at"`
	GroupId     string       `json:"group_id"`
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	SenderId    string       `json:"sender_id"`
	SenderType  SenderType   `json:"sender_type"`
	SourceGuid  string       `json:"source_guid"`
	System      bool         `json:"system"`
	Text        string       `json:"text"`
	UserId      string       `json:"user_id"`
}

type Attachment struct {
	Type    AttachmentType `json:"type"`
	UserIds []string       `json:"user_ids"`
}

type Command interface {
	Name() string
	Matches(msg Message) bool
	Execute(msg Message, c Client) error
}

type CommandBot interface {
	Handler(msg Message) error
}

type CommandBotOptions struct {
	BotId           string
	BotIdFunc       func() string
	AccessToken     string
	AccessTokenFunc func() string
	Logger          *log.Logger
	BaseUrl         string
}

func NewCommandBot(opts CommandBotOptions, cmds ...Command) (CommandBot, error) {
	b := &bot{}

	if len(opts.BotId) > 0 {
		b.botIdFunc = func() string { return opts.BotId }
	} else if opts.BotIdFunc != nil {
		b.botIdFunc = opts.BotIdFunc
	} else {
		return nil, errors.New("bot ID or bot ID func is required")
	}

	if len(opts.AccessToken) > 0 {
		b.accessTokenFunc = func() string { return opts.AccessToken }
	} else if opts.AccessTokenFunc != nil {
		b.accessTokenFunc = opts.AccessTokenFunc
	} else {
		b.accessTokenFunc = func() string { return "" }
	}

	if opts.Logger != nil {
		b.logger = opts.Logger
	} else {
		b.logger = log.New()
	}

	if len(opts.BaseUrl) > 0 {
		b.url = strings.TrimRight(opts.BaseUrl, "/")
	} else {
		b.url = DefaultUrl
	}

	b.commands = cmds
	return b, nil
}

type bot struct {
	botIdFunc       func() string
	accessTokenFunc func() string
	commands        []Command

	logger *log.Logger
	url    string
}

func (b *bot) Handler(msg Message) error {
	msgText := fmt.Sprintf("%v", msg)
	b.logger.Debugf("Received message: %s\n", msgText)

	if msg.SenderType != UserSender {
		b.logger.Debugf("User did not post message, ignoring: %s\n", msgText)
		return nil
	}

	if len(msg.Text) <= 0 {
		b.logger.Debugf("We don't know how to handle empty messages yet: %s\n", msgText)
		return nil
	}

	for _, cmd := range b.commands {
		if cmd.Matches(msg) {
			b.logger.Infof("Found command '%s', executing command on msg: %s\n", cmd.Name(), msgText)

			c := NewClientWithToken(b.botIdFunc(), b.url, b.accessTokenFunc())
			return cmd.Execute(msg, c)
		}
	}

	b.logger.Debugf("No command found for message: %s\n", msgText)
	return nil
}
