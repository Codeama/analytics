// websocket lambda
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/analytics-service/views/processing"
)

// IncomingData represnts data event received
type IncomingData struct {
	ArticleID    string `json:"articleId"`
	ArticleTitle string `json:"articleTitle"`
	PreviousPage string `json:"previousPage"`
	CurrentPage  string `json:"currentPage"`
}

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// todo Send raw events to SNS
	// todo Send tagged event to SNS
	var data IncomingData
	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
		return events.APIGatewayProxyResponse{}, nil
	}

	result, _ := json.Marshal(data)

	fmt.Println("Incoming event:", result)

	var forwardData processing.ViewData
	forwardData.ArticleID = data.ArticleID
	forwardData.ArticleTitle = data.ArticleTitle
	forwardData.PreviousPage = data.PreviousPage
	forwardData.CurrentPage = data.CurrentPage
	forwardData.ConnectionID = request.RequestContext.ConnectionID

	// Process data
	translatedData, _ := processing.TranslateData(forwardData)

	// Marshal data
	processedData, _ := json.Marshal(translatedData)

	fmt.Println("Processed data:", processedData)

	return events.APIGatewayProxyResponse{
		Body:       string(fmt.Sprintf("New lambda: User %s connected!", request.RequestContext.ConnectionID)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
