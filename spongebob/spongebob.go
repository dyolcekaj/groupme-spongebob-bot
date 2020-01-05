package spongebob

import (
	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
	"regexp"
	"unicode"
)

var ptsRxp = *regexp.MustCompile("ok (.*)")

var _ groupme.Command = &PlainTextSarcasm{}
type PlainTextSarcasm struct {

}

func (c *PlainTextSarcasm) Name() string {
	return "PlainTextSarcasm"
}

func (c *PlainTextSarcasm) Matches(text string) bool {
	return ptsRxp.Match([]byte(text))
}

func (c *PlainTextSarcasm) Execute(msg groupme.Message, client *groupme.Client) error {
	return nil
}

var _ groupme.Command = &LastMessageSarcasm{}
type LastMessageSarcasm struct {

}

func (c *LastMessageSarcasm) Name() string {
	return "PlainTextSarcasm"
}

func (c *LastMessageSarcasm) Matches(text string) bool {
	return false
}

func (c *LastMessageSarcasm) Execute(msg groupme.Message, client *groupme.Client) error {
	return nil
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
