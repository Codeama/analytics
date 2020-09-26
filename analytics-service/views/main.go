// websocket lambda
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/codeama/analytics/analytics-service/views/tags"
)

// IncomingData represents data event received
type IncomingData struct {
	ArticleID    string `json:"articleId"`
	ArticleTitle string `json:"articleTitle"`
	PreviousPage string `json:"previousPage"`
	CurrentPage  string `json:"currentPage"`
}

type AttributeValue map[string]*sns.MessageAttributeValue

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

	var forwardData tags.ViewData
	forwardData.ArticleID = data.ArticleID
	forwardData.ArticleTitle = data.ArticleTitle
	forwardData.PreviousPage = data.PreviousPage
	forwardData.CurrentPage = data.CurrentPage
	forwardData.ConnectionID = request.RequestContext.ConnectionID

	// // Process data
	translatedData, _ := tags.TranslateData(forwardData)

	// // Marshal data
	processedData, _ := json.Marshal(translatedData)

	fmt.Println("Processed data:", string(processedData))
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		fmt.Println("NewSession error:", err)
		return events.APIGatewayProxyResponse{}, err
	}
	// tODO extract SNS function
	client := sns.New(sess)
	input := &sns.PublishInput{
		Message: aws.String(string(processedData)),
		MessageAttributes: AttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String("post_views"),
			},
		},
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	}

	_, err = client.Publish(input)
	if err != nil {
		fmt.Println("Publish error:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// fmt.Println(result)

	return events.APIGatewayProxyResponse{
		Body:       string(fmt.Sprintf("New lambda: User %s connected!", request.RequestContext.ConnectionID)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
