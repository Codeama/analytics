package store

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetClient creates a dynamodb client to connect to acsess the datastore
func GetClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TABLE_REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not create a new session: %v", err)
	}

	return dynamodb.New(sess), nil
}
