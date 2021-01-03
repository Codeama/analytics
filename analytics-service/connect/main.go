package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.Headers["Origin"] != os.Getenv("DOMAIN_NAME") {
		fmt.Println("REQUEST HEADERS: ", request.Headers)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, errors.New("Unauthorized")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
