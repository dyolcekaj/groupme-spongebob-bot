package spongebob

import (
	"testing"

	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
)

func TestTranslateText(t *testing.T) {
	testCases := []struct {
		name, text, exp string
	}{
		{
			name: "long string",
			text: "this is a test string",
			exp:  "tHiS iS a TeSt StRiNg",
		},
		{
			name: "one char string",
			text: "a",
			exp:  "a",
		},
		{
			name: "one uc char string",
			text: "A",
			exp:  "a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := translateText(tc.text)
			if act != tc.exp {
				t.Errorf("exp %s, got %s", tc.exp, act)
			}
		})
	}
}

func TestYouKnowWhatSarcasm_Matches(t *testing.T) {
	testCases := []struct {
		name string
		text string
		exp  bool
	}{
		{
			name: "lower case match",
			text: "you know what bud",
			exp:  true,
		},
		{
			name: "mixed case match",
			text: "YoU kNoW wHaT guy",
			exp:  true,
		},
		{
			name: "no match",
			text: "u kno what friend",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := (&YouKnowWhatSarcasm{}).Matches(
				groupme.Message{
					Text: tc.text,
				},
			)
			if act != tc.exp {
				t.Errorf("exp %v, got %v. text: %s", tc.exp, act, tc.text)
			}
		})
	}
}

func TestPlainTextSarcasm_Matches(t *testing.T) {
	testCases := []struct {
		name string
		text string
		exp  bool
	}{
		{"match1", "ok this should match", true},
		{"no match", "this should not match", false},
		{"match2", "Ok this should match", true},
		{"match3", "OK this should match", true},
		{"match4", "oK this should match", true},
		{"no match, only ok", "ok", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := (&CurrentMessageSarcasm{}).Matches(
				groupme.Message{Text: tc.text},
			)

			if act != tc.exp {
				t.Errorf("exp %v, got %v. text: %s", tc.exp, act, tc.text)
			}
		})
	}
}

func TestLastMessageSarcasm_Matches(t *testing.T) {
	testCases := []struct {
		name string
		at   groupme.AttachmentType
		uids []string
		text string
		exp  bool
	}{
		{
			name: "match",
			at:   groupme.MentionType,
			uids: []string{"123"},
			text: "ok @user dude",
			exp:  true,
		},
		{
			name: "no match, bad attachment",
			at:   groupme.LocationType,
			uids: []string{"123"},
			text: "ok @user dude",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := (&LastMessageSarcasm{}).Matches(
				groupme.Message{
					Attachments: []groupme.Attachment{
						{
							Type:    tc.at,
							UserIds: tc.uids,
						},
					},
					Text: tc.text,
				},
			)

			if act != tc.exp {
				t.Errorf("exp %v, got %v", tc.exp, act)
			}
		})
	}
}
