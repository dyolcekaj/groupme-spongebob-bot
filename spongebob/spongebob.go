package spongebob

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
)

var okPrefixRxp = *regexp.MustCompile("^ok (.*)")

var _ groupme.Command = &CurrentMessageSarcasm{}
type CurrentMessageSarcasm struct {

}

func (c *CurrentMessageSarcasm) Name() string {
	return "CurrentMessageSarcasm"
}

func (c *CurrentMessageSarcasm) Matches(msg groupme.Message) bool {
	return okPrefixRxp.Match([]byte(msg.Text)) && len(msg.Attachments) == 0
}

func (c *CurrentMessageSarcasm) Execute(msg groupme.Message, client groupme.Client) error {
	return client.PostBotMessage(translateText(okPrefixRxp.FindStringSubmatch(msg.Text)[1]))
}

var _ groupme.Command = &LastMessageSarcasm{}
type LastMessageSarcasm struct {

}

func (c *LastMessageSarcasm) Name() string {
	return "LastMessageSarcasm"
}

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

	ms, err := client.GetGroupMessages(msg.GroupId, groupme.GroupMessageParams{
		Limit:    100,
	})
	if err != nil {
		return err
	}

	// Messages are sorted in descending time, don't need to sort
	// Can worry about fetching more messages when not found later, 100 should
	// be enough
	for _, m := range ms {
		if m.SenderId == uid {
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
