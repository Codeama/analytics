package store

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/codeama/analytics/analytics-service/post-hits/process"
)

// getClient creates a dynamodb client to connect to acsess the datastore
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
func UpdateTable(data []process.ProcessedEvent) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	for _, article := range data {
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":articleTitle": {
					S: aws.String(article.ArticleTitle),
				},
				":uniqueCount": {
					N: aws.String(strconv.Itoa(article.UniqueViews)),
				},
				":totalCount": {
					N: aws.String(strconv.Itoa(article.TotalViews)),
				},
			},
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"articleId": {
					S: aws.String(article.ArticleID),
				},
			},
			ReturnValues: aws.String("UPDATED_NEW"),

			UpdateExpression: aws.String("ADD uniqueViews :uniqueCount, totalViews :totalCount SET articleTitle = :articleTitle"),
		}

		response, err := client.UpdateItem(input)
		if err != nil {
			return fmt.Errorf("Could not update table item: %v", err)
		}
		fmt.Println("Article updated", response)
	}
	return nil
}
