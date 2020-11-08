// Package publish contains methods that
// process and send events to SNS
package publish

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type attributeValue map[string]*sns.MessageAttributeValue

// AWSSNS is a custom interface for method
// that filters and publishes events to AWS SNS
type AWSSNS interface {
	SendEvent(string) error
}

// Tag represents the tag for the data
// This is used to ToSNS data with the right filter to SNS
type Tag struct {
	Name string
}

func createMessage(eventType string, data string) *sns.PublishInput {
	return &sns.PublishInput{
		Message: aws.String(data),
		MessageAttributes: attributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventType),
			},
		},
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	}
}

// SendEvent implements SSNService interface
// SNS interface (snsiface.SNSAPI) is used here to allow for easy testing and loose coupling
// See here: https://github.com/aws/aws-sdk-go/blob/master/service/sns/snsiface/interface.go
func (tag *Tag) SendEvent(snsClient snsiface.SNSAPI, data string) error {
	var errorMessage = fmt.Errorf("Unable to publish %v", tag)
	switch tag.Name {
	case "homepage_view":
		input := createMessage("homepage_view", data)
		_, err := snsClient.Publish(input)
		if err != nil {
			return errorMessage
		}
	case "post_view":
		input := createMessage("post_view", data)
		_, err := snsClient.Publish(input)
		if err != nil {
			return errorMessage
		}
	case "contact_view":
		input := createMessage("profile_view", data)
		_, err := snsClient.Publish(input)
		if err != nil {
			return errorMessage
		}
	case "about_view":
		input := createMessage("profile_view", data)
		_, err := snsClient.Publish(input)
		if err != nil {
			return errorMessage
		}
	default:
		return fmt.Errorf("Cannot publish data %v", data)
	}
	return nil
}
