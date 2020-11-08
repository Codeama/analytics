// home-hits is a Lambda subscriber function
// to home SQS service
// It processes and counts hits on the homepage
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleQueueMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}
	return nil
}

func main() {
	lambda.Start(handleQueueMessage)
}
