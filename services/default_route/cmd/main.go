// websocket lambda
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/services/default_route"
)

func main() {
	lambda.Start(default_route.HandleRequest)
}
