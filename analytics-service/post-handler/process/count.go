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

func isUnique(previousPage string, currentPage string) bool {
	// UNIQUE: if currentPage is not a / and previousPage is not null
	// HomePage is "/"
	if currentPage != "/" && previousPage != "null" {
		return true
	}
	return false
}

// CountViews totals the number of views for each article
// It filters each article by implementing a set using Go map
// and totals them each one being a view (1)
// Employs a map reduce paradigm
// TODO count unique views
// TODO create a function to extract map values and return an array of items to send to the database
func CountViews(sqsEvent events.SQSEvent) (map[string]ProcessedEvent, error) {
	var data incomingEvent
	var totalViews = make(map[string]int)
	var uniqueViews = make(map[string]int)
	var mappedArt = make(map[string]ProcessedEvent)
	for _, message := range sqsEvent.Records {
		// serialise to Go struct
		if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
			return nil, err
		}
		// checks current article has a view value
		_, hasViews := totalViews[data.ArticleID]
		// _, hasUniqueViews := uniqueViews[data.ArticleID]
		// checks article is already in map (it should be)
		// then updates it (deletes and replaces with article item
		// with the latest view count)
		_, exists := mappedArt[data.ArticleID]
		if hasViews && exists {
			delete(mappedArt, data.ArticleID)

			// process unique views
			unique := isUnique(data.PreviousPage, data.CurrentPage)
			if unique {
				uniqueViews[data.ArticleID]++
			}
			// process all views
			totalViews[data.ArticleID]++

			processed := ProcessedEvent{
				ArticleID:    data.ArticleID,
				ArticleTitle: data.ArticleTitle,
				ConnectionID: data.ConnectionID,
				UniqueViews:  uniqueViews[data.ArticleID],
				TotalViews:   totalViews[data.ArticleID],
			}
			mappedArt[data.ArticleID] = processed
		} else {
			unique := isUnique(data.PreviousPage, data.CurrentPage)
			if unique {
				uniqueViews[data.ArticleID] = 1
			} else {
				uniqueViews[data.ArticleID] = 0
			}
			totalViews[data.ArticleID] = 1
			processed := ProcessedEvent{
				ArticleID:    data.ArticleID,
				ArticleTitle: data.ArticleTitle,
				ConnectionID: data.ConnectionID,
				UniqueViews:  uniqueViews[data.ArticleID],
				TotalViews:   totalViews[data.ArticleID],
			}
			mappedArt[data.ArticleID] = processed
		}

	}

	// return unique items only
	return mappedArt, nil

}
