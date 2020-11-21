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
	session, _ := getSession()

	// publish to SNS
	if err := tag.SendEvent(session, taggedData); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 422}, fmt.Errorf("SendEvent Error: %v", err)
	}

	views := store.ArticleViews{}
	// check incoming data is an article
	if data.ArticleID != "" {
		// read article stats
		views, err = store.GetArticleViews(data.ArticleID)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("ReadTable Error: %v", err)
		}

	}

	stats, err := json.Marshal(views)

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
