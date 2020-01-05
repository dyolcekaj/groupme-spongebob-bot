package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
	"github.com/dyolcekaj/groupme-spongebob-bot/spongebob"
	log "github.com/sirupsen/logrus"
	"net/http"
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

	// Order matters, checked sequentially
	cmds := []groupme.Command{
		&spongebob.LastMessageSarcasm{},
		&spongebob.CurrentMessageSarcasm{},
	}

	bot, err := groupme.NewCommandBot(opts, cmds...)
	if err != nil {
		panic(err)
	}

	a := &App{bot}
	lambda.Start(a.Handler)
}

type App struct {
	GroupMeBot groupme.CommandBot
}

func (a *App) Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Infof("Got request: %v\n", request)
	log.Infof("Body: %s\n", request.Body)
	body := request.Body

	m := groupme.Message{}
	err := json.Unmarshal([]byte(body), &m)
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusInternalServerError,
			Body:              http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	err = a.GroupMeBot.Handler(m)
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusInternalServerError,
			Body:              http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
	}, nil
}