// websocket lambda
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// type ResponseData struct {
// 	ConnectionID string `json:"connectionId"`
// }

// IncomingEvent is the type of event expected
type incomingEvent struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	EventType    string
	ConnectionID string
}

func handleQueueMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	// TODO sort data according to articleId and add/aggregate as single views
	// TODO check previous and currentId to filter for unique views
	// TODO check and retrieve articleId entry in database and update record with total views and unique views

	// receive incoming data event
	var data incomingEvent
	// create a map for counting
	var views = make(map[string]int)
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
			return err
		}

		_, exists := views[data.ArticleID]
		if exists {
			views[data.ArticleID]++
		} else {
			views[data.ArticleID] = 1
		}
	}
	for k, v := range views {
		fmt.Printf("Article ID: %s, Count: %d\n", k, v)
	}
	return nil
}

func main() {
	lambda.Start(handleQueueMessage)
}
