package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dyolcekaj/groupme-spongebob-bot/groupme"
	"log"
	"net/http"
	"os"
	"strings"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var accessToken string

func init() {
	accessToken = os.Getenv("GROUPME_ACCESS_TOKEN")
}

func main() {
	lambda.Start(spongebobify)
}

func spongebobify(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received message: %s\n", req.Body)

	m := &groupme.Message{}
	err := json.Unmarshal([]byte(req.Body), m)
	if err != nil {
		return serverError(err)
	}

	p, err := processMessage(m)
	if err != nil {
		// error indicates we need to ignore the message and return
		log.Printf(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	}

	err = postMessage(m.URL(accessToken), p)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func postMessage(url string, p *groupme.Post) error {
	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Non 200 response code: %v\n", resp)
		return errors.New(msg)
	}

	return nil
}

func processMessage(m *groupme.Message) (*groupme.Post, error) {
	if m.SenderType == groupme.BotSender {
		msg := fmt.Sprintf("Bot user posted message, ignoring: %v\n", m)
		return nil, errors.New(msg)
	}

	if !strings.HasPrefix(m.Text, "SB: ") {
		msg := fmt.Sprintf("Unrelated message, ignoring: %v\n", m)
		return nil, errors.New(msg)
	}
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}
