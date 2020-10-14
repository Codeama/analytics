// websocket lambda
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/analytics-service/views/process"
	"github.com/codeama/analytics/analytics-service/views/publish"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Incoming event:", string(request.Body))
	// TODO Send raw events to SNS
	// todo Send tagged event to SNS
	var data process.IncomingData
	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
		return events.APIGatewayProxyResponse{}, nil
	}

	// validate
	validated, err := process.ValidateData(data, request.RequestContext.ConnectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{}, nil
	}

	// filter
	event := process.FilterData(validated)

	// process and tag
	eventType, taggedData := process.Sort(event)
	fmt.Printf("Tagged data: %v", taggedData)

	// publish to SNS
	publish.SendEvent(eventType, taggedData)

	return events.APIGatewayProxyResponse{
		Body:       string(fmt.Sprintf("New lambda: User %s connected!", request.RequestContext.ConnectionID)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
