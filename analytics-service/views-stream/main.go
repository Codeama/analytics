// Stream lambda is both the websocket backend
// that receives data from the website and does
// the initial data processing for events sent from the blog website
// The aim is to shift data processing logic away from the website
// and avoid load on performance
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
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/codeama/analytics/analytics-service/views-stream/process"
	"github.com/codeama/analytics/analytics-service/views-stream/publish"
	"github.com/codeama/analytics/analytics-service/views-stream/store"
)

type response struct {
	ArticleViews int `json:"uniqueViews"`
}

func getSession() (*sns.SNS, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TOPIC_REGION")),
	})
	if err != nil {
		return &sns.SNS{}, fmt.Errorf("Unable to create session: %v", err)
	}
	session := sns.New(sess)
	return session, nil
}

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Incoming event:", string(request.Body))
	var data process.IncomingData
	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("Error: %v", err)
	}

	// validate
	validated, err := process.ValidateData(data, request.RequestContext.ConnectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("Error: %v", err)
	}

	// filter
	event := process.FilterData(validated)

	// process and tag
	eventType, taggedData, _ := process.Sort(event)
	fmt.Printf("Tagged data: %v", taggedData)

	// convert eventType to a Tag struct
	tag := &publish.Tag{Name: eventType}

	// get config
	snsSession, _ := getSession()

	// publish to SNS
	if err := tag.SendEvent(snsSession, taggedData); err != nil {
		fmt.Println("Send event failed:", err)
		return events.APIGatewayProxyResponse{StatusCode: 422}, fmt.Errorf("SendEvent Error: %v", err)
	}

	connectionSess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TOPIC_REGION")), //TODO refactor session logic
	})

	if err != nil {
		fmt.Println("Connection session failed.")
		// return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("Session Error: %v", err)
	}

	views := store.ArticleViews{}
	// check incoming data is an article
	if data.ArticleID != "" {
		// read article stats
		fmt.Println("Attempting to send to Connected client...")
		views, err = store.GetArticleViews(data.ArticleID)
		if err != nil {
			fmt.Println("Couldn't get article views from data store.")
			// return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("ReadTable Error: %v", err)
		}
	}

	viewsStats, err := json.Marshal(views)

	apigw := apigatewaymanagementapi.New(connectionSess, &aws.Config{
		Endpoint: aws.String("https://bplqdxpti2.execute-api.eu-west-2.amazonaws.com/test"),
	})
	_, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(request.RequestContext.ConnectionID),
		Data:         viewsStats,
	})
	if err != nil {
		fmt.Println("Post to connection failed: ", err)
		// return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("PostToConnection Error: %v", err)
	}
	// fmt.Println("PostToCONNECTION: ", output)

	stats, err := json.Marshal(request.RequestContext.ConnectionID)
	if err != nil {
		fmt.Println("Could not marshall request ID")
	}

	return events.APIGatewayProxyResponse{
		// Mehh when you finish a project and realise you don't need websockets >.<
		// But I learned a lot...
		Body:       string(stats),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
