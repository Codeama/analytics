package store

import (
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/codeama/analytics/analytics-service/views-stream/process"
	"github.com/google/uuid"
)

// UpdateReferrerTable updates the table with the given value from SQS event
func UpdateReferrerTable(data process.AnalyticsData) error {
	client, err := GetClient()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	referrer := data.Referrer
	isCurrentDomain, _ := regexp.MatchString(os.Getenv("DOMAIN_NAME"), referrer)

	// Check referrer is external
	if !isCurrentDomain {
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":referrer": {
					S: aws.String(referrer),
				},
				":currentPage": {
					S: aws.String(data.CurrentPage),
				},
				":connectionID": {
					S: aws.String(data.ConnectionID),
				},
			},
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(uuid.New().String()),
				},
			},
			ReturnValues: aws.String("UPDATED_NEW"),

			UpdateExpression: aws.String("SET connectionID = :connectionID, referrerURL = :referrer, pageViewed = :currentPage"),
		}

		response, err := client.UpdateItem(input)
		if err != nil {
			return fmt.Errorf("Could not update table item: %v", err)
		}
		fmt.Println("New referrer URL added to table", response)
	}

	return nil
}
