package process

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// IncomingEvent is the type of event expected
type incomingEvent struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	EventType    string
	ConnectionID string
}

// ProcessedEvent is the result returned after counting
type ProcessedEvent struct {
	ArticleID    string
	ArticleTitle string
	ConnectionID string
	UniqueViews  int
	TotalViews   int // includes a count of non-unique views
}

// CountViews totals the number of views for each article
// It filters for each unique article by implementing a set using Go map
// and totals them each one being a view (1)
// Employs a map reduce paradigm
// TODO count unique views
// TODO create a function to extract map values and return an array of items to send to the database
func CountViews(sqsEvent events.SQSEvent) (map[string]ProcessedEvent, error) {
	var data incomingEvent
	var views = make(map[string]int)
	var mappedArt = make(map[string]ProcessedEvent)
	for _, message := range sqsEvent.Records {
		// serialise to Go struct
		if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
			return nil, err
		}
		_, exists := views[data.ArticleID]
		if exists {
			_, article := mappedArt[data.ArticleID]
			if article {
				delete(mappedArt, data.ArticleID)
			}
			views[data.ArticleID]++
			processed := ProcessedEvent{
				ArticleID:    data.ArticleID,
				ArticleTitle: data.ArticleTitle,
				ConnectionID: data.ConnectionID,
				UniqueViews:  views[data.ArticleID],
				TotalViews:   views[data.ArticleID],
			}
			mappedArt[data.ArticleID] = processed
		} else {
			views[data.ArticleID] = 1
			processed := ProcessedEvent{
				ArticleID:    data.ArticleID,
				ArticleTitle: data.ArticleTitle,
				ConnectionID: data.ConnectionID,
				UniqueViews:  views[data.ArticleID],
				TotalViews:   views[data.ArticleID],
			}
			mappedArt[data.ArticleID] = processed
		}

	}

	// return unique items only
	return mappedArt, nil

}
