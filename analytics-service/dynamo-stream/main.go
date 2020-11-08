// This Lambda reads data from DynamoDB stream and writes to
// another table (for publishing views)
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/analytics-service/dynamo-stream/stream"
)

func handleQueueMessage(ctx context.Context, streamData events.DynamoDBEvent) error {

	if err := stream.UpdateTable(streamData); err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	return nil
}

func main() {
	lambda.Start(handleQueueMessage)
}
