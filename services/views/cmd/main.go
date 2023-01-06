// Stream lambda is both the websocket backend
// that receives data from the website and does
// the initial data processing for events sent from the blog website
// The aim is to shift data processing logic away from the website
// and avoid load on performance
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codeama/analytics/services/views"
)

func main() {
	lambda.Start(views.HandleRequest)
}
