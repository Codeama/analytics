// Package store reads total article views from
// DynamoDB article table (reader) and publishes
// to the websocket
package store

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// ArticleViews represents stats read from table
type ArticleViews struct {
	UniqueViews int `json:"uniqueViews"`
}

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

// GetArticleViews queries the article table and returns
// the value of total views for an article
func GetArticleViews(articleID string) (ArticleViews, error) {
	client, err := getClient()
	if err != nil {
		return ArticleViews{}, fmt.Errorf("Error: %v", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"articleId": {
				S: aws.String(articleID),
			},
		},
	}

	result, err := client.GetItem(input)
	if result.Item == nil {
		return ArticleViews{}, nil
	}
	if err != nil {
		return ArticleViews{}, fmt.Errorf("%v", err)
	}

	item := ArticleViews{}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
		return ArticleViews{}, fmt.Errorf("Failed to unmarshal table result: %v", err)
	}

	return ArticleViews{item.UniqueViews}, nil
}
