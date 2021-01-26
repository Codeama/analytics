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
	ArticleViews int `json:"totalViews"`
}

func getSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to create session: %v", err)
	}
	return sess, nil
}

func getSNSService(session *session.Session) *sns.SNS {
	return sns.New(session)
}

func sendStats(data store.ArticleViews, session *session.Session, connection string) {

	apigw := apigatewaymanagementapi.New(session, &aws.Config{
		Endpoint: aws.String(os.Getenv("CONNECTION_URL")),
	})

	viewsStats, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Marshal data failed: ", err)
	}

	_, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connection),
		Data:         viewsStats,
	})
	if err != nil {
		fmt.Println("Post to connection failed: ", err)
	}
}

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Incoming event:", string(request.Body))
	var data process.IncomingData

	// If error, return as data is useless
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
	fmt.Printf("Tag data: %v", taggedData)

	// convert eventType to a Tag struct
	tag := &publish.Tag{Name: eventType}

	// get session
	session, err := getSession()
	if session != nil {
		// get SNS service
		snsService := getSNSService(session)
		// publish to SNS
		if err := tag.SendEvent(snsService, taggedData); err != nil {
			fmt.Println("Send event failed:", err)
		}

		// If stream data is article, retrieve stats and post to connection
		views := store.ArticleViews{}
		if data.ArticleID != "" {
			views, err = store.GetArticleViews(data.ArticleID)
			if err != nil {
				fmt.Println("Couldn't get article views from data store.", err)
			}
			// Send stats
			sendStats(views, session, request.RequestContext.ConnectionID)
		}

	}

	// Store external referrer URL
	referrer := process.AnalyticsData{
		ArticleID:    data.ArticleID,
		ArticleTitle: data.ArticleTitle,
		CurrentPage:  data.CurrentPage,
		PreviousPage: data.PreviousPage,
		ConnectionID: request.RequestContext.ConnectionID,
		Refreshed:    data.Refreshed,
		Referrer:     data.Referrer,
	}
	if err := store.UpdateReferrerTable(referrer); err != nil {
		fmt.Println("Can't update referrer: ", err)
	}

	return events.APIGatewayProxyResponse{
		// No need to send body data as this is sent via websocket connection
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
