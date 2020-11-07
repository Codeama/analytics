// Post lambda that is subscribed to a PostSQS
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/analytics-service/post-hits/process"
	"github.com/codeama/analytics/analytics-service/post-hits/store"
)

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
	count, err := process.CountViews(sqsEvent)
	if err != nil {
		return fmt.Errorf("Cannot process event: %v", err)
	}

	articles := process.GetCountedPosts(count)
	fmt.Println("ARTICLES: ", articles)

	if err := store.UpdateTable(articles); err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	return nil
}

func main() {
	lambda.Start(handleQueueMessage)
}
