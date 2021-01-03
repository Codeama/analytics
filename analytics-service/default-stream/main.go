// websocket lambda
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		Body:       string(fmt.Sprintln("Hello!!! Nothing to see here!")),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
