package default_route

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func HandleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		Body:       string(fmt.Sprintln("Hello!!! Nothing to see here!")),
		StatusCode: 200,
	}, nil
}
