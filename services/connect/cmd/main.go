package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/services/connect"
)

func main() {
	lambda.Start(connect.HandleRequest)
}
