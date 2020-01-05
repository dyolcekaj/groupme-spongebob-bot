package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"

	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
	"github.com/dyolcekaj/groupme-spongebob-bot/spongebob"
)

var (
	botName = os.Getenv("GROPUME_BOT_NAME")
	accessToken = os.Getenv("GROUPME_ACCESS_TOKEN")
)

func main() {
	opts := groupme.CommandBotOptions{
		AccessToken: accessToken,
	}

	// Order matters, checked sequentially
	cmds := []groupme.Command{
		&spongebob.LastMessageSarcasm{},
		&spongebob.CurrentMessageSarcasm{},
	}

	bot, err := groupme.NewCommandBot(botName, opts, cmds...)
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
	body := request.Body

	m := groupme.Message{}
	err := json.Unmarshal([]byte(body), &m)
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	err = a.GroupMeBot.Handler(m)
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
