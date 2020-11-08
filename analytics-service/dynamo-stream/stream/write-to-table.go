// Package stream contains functions
// that process iincoming dynamodb stream events
package stream

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// getClient creates a dynamodbstream client to connect to acsess the datastore
func getClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TABLE_REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not create a new session: %v", err)
	}

	return dynamodb.New(sess), nil
}

// UpdateTable updates the table with the given slice value
func UpdateTable(streamData events.DynamoDBEvent) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	for _, record := range streamData.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":articleTitle": {
					S: aws.String(record.Change.NewImage["articleTitle"].String()),
				},
				":uniqueCount": {
					N: aws.String(record.Change.NewImage["uniqueViews"].Number()),
				},
				":totalCount": {
					N: aws.String(record.Change.NewImage["totalViews"].Number()),
				},
			},
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"articleId": {
					S: aws.String(record.Change.NewImage["articleId"].String()),
				},
			},
			ReturnValues: aws.String("UPDATED_NEW"),

			UpdateExpression: aws.String("SET uniqueViews = :uniqueCount, totalViews = :totalCount, articleTitle = :articleTitle"),
		}

		response, err := client.UpdateItem(input)
		if err != nil {
			return fmt.Errorf("Could not update table item: %v", err)
		}
		fmt.Println("Article updated", response)

	}
	return nil
}
