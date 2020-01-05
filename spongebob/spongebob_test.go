package spongebob

import (
	"github.com/dyolcekaj/groupme-spongebob-bot/assertions"
	"testing"
)

func TestTranslateText(t *testing.T) {
	ret := translateText("this is a test string")
	assertions.Equals(t, "tHiS iS a TeSt StRiNg", ret)
}