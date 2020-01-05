package spongebob

import (
	"github.com/dyolcekaj/groupme-spongebob-bot/assertions"
	"testing"
)

func TestTranslateText(t *testing.T) {
	ret := translateText("this is a test string")
	assertions.Equals(t, "tHiS iS a TeSt StRiNg", ret)
}

func TestPlainTextSarcasm_Matches(t *testing.T) {
	pts := &PlainTextSarcasm{}
	assertions.Assert(t, pts.Matches("ok this should match"), "Should have matched")
	assertions.Assert(t, !pts.Matches("this should not match"), "Should not have matched")
}

func TestLastMessageSarcasm_Matches(t *testing.T) {

}