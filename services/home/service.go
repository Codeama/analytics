package home

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func HandleQueueMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	var viewData IncomingEvent
	for _, message := range sqsEvent.Records {
		// serialise to Go struct
		if err := json.Unmarshal([]byte(message.Body), &viewData); err != nil {
			return fmt.Errorf("could not deserialise data: %v", err)
		}

		if err := UpdateTable(viewData); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
