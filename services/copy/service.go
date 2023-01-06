package copy

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func HandleQueueMessage(ctx context.Context, streamData events.DynamoDBEvent) error {

	if err := UpdateTable(streamData); err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	return nil
}
