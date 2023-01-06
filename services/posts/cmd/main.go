// Post lambda that is subscribed to a PostSQS
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/services/posts"
)

func main() {
	lambda.Start(posts.HandleQueueMessage)
}
