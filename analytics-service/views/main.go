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

// IncomingData represents data event received
type IncomingData struct {
	ArticleID    string `json:"articleId,omitempty"`
	ArticleTitle string `json:"articleTitle,omitempty"`
	PreviousPage string `json:"previousPage"`
	CurrentPage  string `json:"currentPage"`
}

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// todo Send raw events to SNS
	// todo Send tagged event to SNS
	var data IncomingData
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{}, nil
	}

	result, _ := json.Marshal(data)

	fmt.Println("Incoming event:", string(result))

	// turn to a func that returns unknown interface
	var received process.ReceivedData
	received.ArticleID = data.ArticleID
	received.ArticleTitle = data.ArticleTitle
	received.PreviousPage = data.PreviousPage
	received.CurrentPage = data.CurrentPage
	received.ConnectionID = request.RequestContext.ConnectionID

	filtered := process.FilterData(received)

	eventType, taggedData := process.Sort(filtered)

	fmt.Printf("Tagged data: %v", taggedData)

	// processedData, _ := json.Marshal(taggedData)

	// fmt.Println("Processed data:", string(processedData))
	// publish to SNS
	// TODO check args are not empty strings???
	publish.SendEvent(eventType, taggedData)

	return events.APIGatewayProxyResponse{
		Body:       string(fmt.Sprintf("New lambda: User %s connected!", request.RequestContext.ConnectionID)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
