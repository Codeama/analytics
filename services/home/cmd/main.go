// home-hits is Lambda that subscribes
// to and reads messages from 'home' SQS service
// It then processes and counts hits on the homepage
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/services/home"
)

func main() {
	lambda.Start(home.HandleQueueMessage)
}
