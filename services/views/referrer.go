package views

import (
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// IsDomain checks whether referrer value is same as domain origin
func IsDomain(ref string) bool {
	isCurrentDomain, _ := regexp.MatchString(os.Getenv("DOMAIN_NAME"), ref)
	return isCurrentDomain
}

// UpdateReferrerTable only updates the table if referrer is external
// and not null
func UpdateReferrerTable(data AnalyticsData) error {
	client, err := GetClient()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	isCurrentDomain := IsDomain(data.Referrer)

	// Only store referrer if it is external
	if !isCurrentDomain && data.Referrer != "" {
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":referrer": {
					S: aws.String(data.Referrer),
				},
				":currentPage": {
					S: aws.String(data.CurrentPage),
				},
				":connectionID": {
					S: aws.String(data.ConnectionID),
				},
			},
			TableName: aws.String(os.Getenv("REFERRER_TABLE_NAME")),
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
