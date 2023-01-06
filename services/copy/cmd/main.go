// This Lambda reads data from DynamoDB stream and writes to
// another table (for publishing views)
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/services/copy"
)

func main() {
	lambda.Start(copy.HandleQueueMessage)
}
