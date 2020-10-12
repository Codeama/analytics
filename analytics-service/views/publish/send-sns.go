// Package publish adds appropriate event type and sends to SNS
package publish

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AttributeValue map[string]*sns.MessageAttributeValue

func publish(eventType string, data string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	client := sns.New(sess)
	input := &sns.PublishInput{
		Message: aws.String(data),
		MessageAttributes: AttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventType),
			},
		},
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	}

	_, err = client.Publish(input)
	if err != nil {
		fmt.Println("Publish error:", err)
		return
	}
}

// SendEvent publishes events to SNS
func SendEvent(tagName string, data string) {
	switch tagName {
	case "post_view":
		publish("post_view", data)
	case "contact_view":
		publish("profile_view", data)
	case "about_view":
		publish("profile_view", data)
	case "homepage_view":
		// publish("homepage_view", data)
		fmt.Println("HOMEPAGE VIEW", data)
	default:
		fmt.Println("What is this data?", data)
		// publish("raw_data", data)
	}

}
