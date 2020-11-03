// Homepage view queue subscriber
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleQueueMessage(ctx context.Context) error {
	fmt.Println("Hello, DynamoDB streams!")
	return nil
}

func main() {
	lambda.Start(handleQueueMessage)
}
