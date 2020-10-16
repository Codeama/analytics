// Views lambda is both the websocket connection backend
// that receives data from the website and
// the initial data processing 'frontend'
// for events sent from the blog website
// The aim is to shift data processing logic away from the website
// to avoid load on performance
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
	eventType, taggedData, _ := process.Sort(event)
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
