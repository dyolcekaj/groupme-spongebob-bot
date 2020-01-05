package spongebob

import (
	"github.com/dyolcekaj/groupme-spongebob-bot/assertions"
	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
	"testing"
)

func TestTranslateText(t *testing.T) {
	ret := translateText("this is a test string")
	assertions.Equals(t, "tHiS iS a TeSt StRiNg", ret)
}

func TestPlainTextSarcasm_Matches(t *testing.T) {
	pts := &CurrentMessageSarcasm{}

	assertions.Assert(t, pts.Matches(groupme.Message{Text: "ok this should match"}), "Should have matched")
	assertions.Assert(t, !pts.Matches(groupme.Message{Text: "this should not match"}), "Should not have matched")
}

func TestLastMessageSarcasm_Matches(t *testing.T) {
	lms := &LastMessageSarcasm{}

	assertions.Assert(t, lms.Matches(groupme.Message{
		Attachments: []groupme.Attachment{
			{
				Type:    groupme.MentionType,
				UserIds: []string{"123"},
			},
		},
		Text:        "ok @user dude",
	}), "should match")

	assertions.Assert(t, !lms.Matches(groupme.Message{
		Attachments: []groupme.Attachment{
			{
				Type:    groupme.LocationType,
				UserIds: []string{"123"},
			},
		},
		Text:        "ok @user dude",
	}), "should not match")
}