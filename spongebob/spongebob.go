package spongebob

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
)

var okPrefixRxp = *regexp.MustCompile("^[Oo][Kk] (.*)")
var youKnowWhatPrefix = "you know what"

var _ groupme.Command = &YouKnowWhatSarcasm{}

// YouKnowWhatSarcasm responds sarcastically if a sentence is started
// with the phrase "you know what"
type YouKnowWhatSarcasm struct{}

// Name returns name of this comnmand
func (c *YouKnowWhatSarcasm) Name() string {
	return "CurrentMessageSarcasm"
}

// Matches a GroupMe message
func (c *YouKnowWhatSarcasm) Matches(msg groupme.Message) bool {
	return len(msg.Text) > 0 && strings.HasPrefix(strings.ToLower(msg.Text), youKnowWhatPrefix)
}

// Execute a response to a matched GroupMe message
func (c *YouKnowWhatSarcasm) Execute(msg groupme.Message, client groupme.Client) error {
	return client.PostBotMessage(translateText(msg.Text))
}

var _ groupme.Command = &CurrentMessageSarcasm{}

// CurrentMessageSarcasm responds sarcastically if a sentence starts with
// "ok " but has no user mentions
type CurrentMessageSarcasm struct {
}

// Name returns name of this comnmand
func (c *CurrentMessageSarcasm) Name() string {
	return "CurrentMessageSarcasm"
}

// Matches a GroupMe message
func (c *CurrentMessageSarcasm) Matches(msg groupme.Message) bool {
	return okPrefixRxp.Match([]byte(msg.Text)) && len(msg.Attachments) == 0
}

// Execute a response to a matched GroupMe message
func (c *CurrentMessageSarcasm) Execute(msg groupme.Message, client groupme.Client) error {
	return client.PostBotMessage(translateText(okPrefixRxp.FindStringSubmatch(msg.Text)[1]))
}

var _ groupme.Command = &LastMessageSarcasm{}

// LastMessageSarcasm sarcastically repeats back whatever the mentioned user said last
// if the current message starts with "ok @mention"
type LastMessageSarcasm struct {
}

// Name returns name of this comnmand
func (c *LastMessageSarcasm) Name() string {
	return "LastMessageSarcasm"
}

// Matches a GroupMe message
func (c *LastMessageSarcasm) Matches(msg groupme.Message) bool {
	if !okPrefixRxp.Match([]byte(msg.Text)) || len(msg.Attachments) == 0 {
		return false
	}

	for _, a := range msg.Attachments {
		if a.Type == groupme.MentionType {
			// has at least one mention
			return true
		}
	}

	return false
}

// Execute a response to a matched GroupMe message
func (c *LastMessageSarcasm) Execute(msg groupme.Message, client groupme.Client) error {
	var uid string
	for _, a := range msg.Attachments {
		if a.Type == groupme.MentionType {
			uid = a.UserIds[0]
			break
		}
	}

	if len(uid) == 0 {
		return fmt.Errorf("no user mentioned in message, can't search: %s", msg.Text)
	}

	ms, err := client.GetGroupMessages(
		groupme.GroupMessageParams{
			Limit: 100,
		},
	)
	if err != nil {
		return err
	}

	// Messages are sorted in descending time, don't need to sort
	// Can worry about fetching more messages when not found later, 100 should
	// be enough. Only respond to user messages with no attachments as a
	// quick and dirty default
	for _, m := range ms {
		if m.SenderType == groupme.UserSender && m.SenderID == uid && len(m.Attachments) == 0 {
			return client.PostBotMessage(translateText(m.Text))
		}
	}

	return fmt.Errorf("no messages found for uid: %s", uid)
}

func translateText(text string) string {
	var ret []rune
	i := false
	for _, r := range text {
		if i {
			ret = append(ret, unicode.ToUpper(r))
		} else {
			ret = append(ret, unicode.ToLower(r))
		}

		if unicode.IsLetter(r) {
			i = !i
		}
	}

	return string(ret)
}
