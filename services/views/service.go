package views

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/sns"
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

func sendStats(data ArticleViews, session *session.Session, connection string) {

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

func HandleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Incoming event:", string(request.Body))
	var data IncomingData

	// If error, return as data is useless
	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("Error: %v", err)
	}

	// validate
	validated, err := ValidateData(data, request.RequestContext.ConnectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("Error: %v", err)
	}

	// filter
	event := FilterData(validated)

	// process and tag
	eventType, taggedData, _ := Sort(event)
	fmt.Printf("Tag data: %v", taggedData)

	// convert eventType to a Tag struct
	tag := &Tag{Name: eventType}

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
		views := ArticleViews{}
		if data.ArticleID != "" {
			views, err = GetArticleViews(data.ArticleID)
			if err != nil {
				fmt.Println("Couldn't get article views from data store.", err)
			}
			// Send stats
			sendStats(views, session, request.RequestContext.ConnectionID)
		}

	}

	// Store external referrer URL
	referrer := AnalyticsData{
		ArticleID:    data.ArticleID,
		ArticleTitle: data.ArticleTitle,
		CurrentPage:  data.CurrentPage,
		PreviousPage: data.PreviousPage,
		ConnectionID: request.RequestContext.ConnectionID,
		Refreshed:    data.Refreshed,
		Referrer:     data.Referrer,
	}
	if err := UpdateReferrerTable(referrer); err != nil {
		fmt.Println("Can't update referrer: ", err)
	}

	return events.APIGatewayProxyResponse{
		// No need to send body data as this is sent via websocket connection
		StatusCode: 200,
	}, nil
}
