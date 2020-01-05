package groupme

import (
	"errors"
	"fmt"
	"github.com/dyolcekaj/groupme-spongebob-bot/groupme/internal"
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

func (m *Message) String() string {
	return fmt.Sprintf(
		"[createAt: %d, groupId: %s, id: %s, name: %s, senderId: %s, senderType: %s, text: %s, userId: %s]",
		m.CreatedAt, m.GroupId, m.Id, m.Name, m.SenderId, m.SenderType, m.Text, m.UserId,
	)
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
	AccessToken   string
	Logger        *log.Logger
	BaseUrl       string
}

func NewCommandBot(name string, opts CommandBotOptions, cmds ...Command) (CommandBot, error) {
	b := &bot{name: name}

	if len(opts.AccessToken) > 0 {
		b.accessToken = opts.AccessToken
	} else {
		return nil, errors.New("access token is a required opt")
	}

	if len(cmds) == 0 {
		return nil, errors.New("no commands provided")
	}
	b.commands = cmds

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

	b.cache = internal.NewCache()
	b.cacheClient = NewClient("", b.url, b.accessToken)

	err := b.loadCache()
	if err != nil {
		return nil, fmt.Errorf("error loading bot ids for botname '%s': %w", name, err)
	}

	return b, nil
}

type bot struct {
	name        string
	accessToken string
	commands    []Command

	logger *log.Logger
	url    string

	cache       internal.BotIdCache
	cacheClient Client
}

func (b *bot) Handler(msg Message) error {
	msgText := msg.String()
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

			botId, ok := b.findBotId(msg.GroupId)
			if !ok {
				err := fmt.Errorf(
					"No bot id found for command '%s' and groupd id '%s' on msg: %s\n",
					cmd.Name(), msg.GroupId, msgText,
				)
				b.logger.Error(err)
				return err
			}

			c := NewClient(botId, b.url, b.accessToken)
			return cmd.Execute(msg, c)
		}
	}

	b.logger.Debugf("No command found for message: %s\n", msgText)
	return nil
}

func (b *bot) findBotId(groupId string) (string, bool) {
	if botId, ok := b.cache.Get(groupId); ok {
		return botId, ok
	}

	// try, try again
	err := b.loadCache()
	if err != nil {
		b.logger.Errorf("error reloading cache: %v", err)
	}
	return b.cache.Get(groupId)
}

func (b *bot) loadCache() error {
	bots, err := b.cacheClient.ListBots()
	if err != nil {
		return err
	}

	b.cache.Clear()
	for _, bot := range bots {
		b.cache.Set(bot.GroupId, bot.BotId)
	}
	return nil
}
