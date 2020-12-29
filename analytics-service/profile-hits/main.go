// profile-hits is Lambda that subscribes
// to and reads messages from 'home' SQS service
// It then processes and counts hits on the homepage
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/analytics-service/profile-hits/store"
)

func handleQueueMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	var viewData store.IncomingEvent
	for _, message := range sqsEvent.Records {
		// serialise to Go struct
		if err := json.Unmarshal([]byte(message.Body), &viewData); err != nil {
			return fmt.Errorf("Could not deserialise data: %v", err)
		}

		if err := store.UpdateTable(viewData); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func main() {
	lambda.Start(handleQueueMessage)
}
