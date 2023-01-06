package home

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// IncomingEvent represents struct for the expected SQS event data
type IncomingEvent struct {
	ConnectionID string
	CurrentPage  string
	PreviousPage string
	Refreshed    bool
	EventType    string
	Referrer     string
}

// getClient creates a dynamodb client to connect DynamoDB to acsess the datastore
func getClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TABLE_REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not create a new session: %v", err)
	}

	return dynamodb.New(sess), nil
}

// UpdateTable updates the table with the given value from SQS event
func UpdateTable(data IncomingEvent) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	if data.EventType == "homepage_view" && data.ConnectionID != "" {
		var uniqueCount int
		if !data.Refreshed {
			uniqueCount = 1
		}

		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":totalCount": {
					N: aws.String(strconv.Itoa(1)),
				},
				":uniqueCount": {
					N: aws.String(strconv.Itoa(uniqueCount)),
				},
			},
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"pageName": {
					S: aws.String("Home_Page"),
				},
			},
			ReturnValues: aws.String("UPDATED_NEW"),

			UpdateExpression: aws.String("ADD uniqueViews :uniqueCount, totalViews :totalCount"),
		}

		response, err := client.UpdateItem(input)
		if err != nil {
			return fmt.Errorf("Could not update table item: %v", err)
		}
		fmt.Println("Home_Page item updated", response)
	}

	return nil
}
