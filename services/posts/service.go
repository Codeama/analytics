package posts

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func HandleQueueMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	count, err := CountViews(sqsEvent)
	if err != nil {
		return fmt.Errorf("Cannot process event: %v", err)
	}

	articles := GetCountedPosts(count)
	fmt.Println("ARTICLES: ", articles)

	if err := UpdateTable(articles); err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	return nil
}
