package groupme

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dyolcekaj/groupme-spongebob-bot/groupme/internal"
	log "github.com/sirupsen/logrus"
)

// SenderType is type of message sender
type SenderType string

// AttachmentType is type of attachment on message
type AttachmentType string

// Defined SenderTypes and AttachmentTypes
const (
	UserSender SenderType = "user"
	BotSender  SenderType = "bot"

	ImageType    AttachmentType = "image"
	LocationType AttachmentType = "location"
	MentionType  AttachmentType = "mentions"
	EmojiType    AttachmentType = "emoji"
	SplitType    AttachmentType = "split"

	DefaultURL = "https://api.groupme.com/v3"
)

// Message sent during callback from GroupMe
type Message struct {
	Attachments []Attachment `json:"attachments"`
	CreatedAt   int          `json:"created_at"`
	GroupID     string       `json:"group_id"`
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	SenderID    string       `json:"sender_id"`
	SenderType  SenderType   `json:"sender_type"`
	SourceGUID  string       `json:"source_guid"`
	System      bool         `json:"system"`
	Text        string       `json:"text"`
	UserID      string       `json:"user_id"`
}

func (m *Message) String() string {
	return fmt.Sprintf(
		"[createAt: %d, groupId: %s, id: %s, name: %s, senderId: %s, senderType: %s, text: %s, userId: %s]",
		m.CreatedAt, m.GroupID, m.ID, m.Name, m.SenderID, m.SenderType, m.Text, m.UserID,
	)
}

// Attachment that may/may not be on message
type Attachment struct {
	Type    AttachmentType `json:"type"`
	UserIds []string       `json:"user_ids"`
}

// Command to match against and act upon GroupMe Messages
type Command interface {
	Name() string
	Matches(msg Message) bool
	Execute(msg Message, c Client) error
}

// CommandBot to handle Commands
type CommandBot interface {
	Handler(msg Message) error
}

// CommandBotOptions for configuring the bot
type CommandBotOptions struct {
	AccessToken string
	Logger      *log.Logger
	BaseURL     string
}

// NewCommandBot returns a configured CommandBot
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

	if len(opts.BaseURL) > 0 {
		b.url = strings.TrimRight(opts.BaseURL, "/")
	} else {
		b.url = DefaultURL
	}

	b.cache = internal.NewCache()
	b.cacheClient = internal.NewClient(b.url, b.accessToken)

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

	cache       internal.BotIDCache
	cacheClient internal.Client
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

	bot, ok := b.findBot(msg.GroupID)
	if !ok {
		err := fmt.Errorf(
			"no bot id found for groupd id '%s' on msg: %s",
			msg.GroupID, msgText,
		)
		b.logger.Error(err)
		return err
	}

	for _, cmd := range b.commands {
		if cmd.Matches(msg) {
			b.logger.Infof("Found command '%s', executing command on msg: %s\n", cmd.Name(), msgText)

			c := newClient(bot.BotID, bot.GroupID, b.url, b.accessToken)
			return cmd.Execute(msg, c)
		}
	}

	b.logger.Debugf("No command found for message: %s\n", msgText)
	return nil
}

func (b *bot) findBot(groupID string) (internal.Bot, bool) {
	if bot, ok := b.cache.Get(groupID); ok {
		return bot, ok
	}

	// try, try again
	err := b.loadCache()
	if err != nil {
		b.logger.Errorf("error reloading cache: %v", err)
	}
	return b.cache.Get(groupID)
}

func (b *bot) loadCache() error {
	bots, err := b.cacheClient.ListBots()
	if err != nil {
		return err
	}

	b.cache.Clear()
	for _, bot := range bots {
		b.cache.Set(bot.GroupID, bot)
	}
	return nil
}
