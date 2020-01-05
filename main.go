package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
	"github.com/dyolcekaj/groupme-spongebob-bot/spongebob"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	l := log.New()
	l.SetLevel(log.DebugLevel)

	opts := groupme.CommandBotOptions{
		BotIdFunc:       func() string { return os.Getenv("GROUPME_BOT_ID") },
		AccessTokenFunc: func() string { return os.Getenv("GROUPME_ACCESS_TOKEN") },
		Logger:          l,
		BaseUrl:         groupme.DefaultUrl,
	}

	cmds := []groupme.Command{
		&spongebob.PlainTextSarcasm{},
		&spongebob.LastMessageSarcasm{},
	}

	bot, err := groupme.NewCommandBot(opts, cmds...)
	if err != nil {
		panic(err)
	}

	lambda.Start(bot.Handler())
}
