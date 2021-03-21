package process

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// IncomingEvent is the type of event expected
type incomingEvent struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	Refreshed    bool
	EventType    string
	ConnectionID string
	Referrer     string
}

// ProcessedEvent is the result returned after counting
type ProcessedEvent struct {
	ArticleID    string
	ArticleTitle string
	UniqueViews  int
	TotalViews   int // sum total of all views unique or not
}

// CountViews totals the number of views for each article
// It processes the sum of each article viewed by implementing a set using Go map
// and does parallel data processing
func CountViews(sqsEvent events.SQSEvent) (map[string]ProcessedEvent, error) {
	var data incomingEvent
	var totalViews = make(map[string]int)
	var uniqueViews = make(map[string]int)
	var mappedArt = make(map[string]ProcessedEvent)
	for _, message := range sqsEvent.Records {
		// serialise to Go struct
		if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
			return nil, fmt.Errorf("could not deserialise data: %v", err)
		}
		// checks current article has a view value
		_, hasViews := totalViews[data.ArticleID]
		// checks article is already in map (it should be; see else statement below that runs at least once for all incoming data)
		// then updates it (by deleting existing one and replacing it with latest stats)
		_, exists := mappedArt[data.ArticleID]
		if hasViews && exists {
			delete(mappedArt, data.ArticleID)

			// process unique views
			if !data.Refreshed {
				uniqueViews[data.ArticleID]++
			}
			// process all views
			totalViews[data.ArticleID]++

			processed := ProcessedEvent{
				ArticleID:    data.ArticleID,
				ArticleTitle: data.ArticleTitle,
				UniqueViews:  uniqueViews[data.ArticleID],
				TotalViews:   totalViews[data.ArticleID],
			}
			// update article values
			mappedArt[data.ArticleID] = processed
		} else {
			// process unique views
			if !data.Refreshed {
				uniqueViews[data.ArticleID] = 1
			} else {
				uniqueViews[data.ArticleID] = 0
			}
			// process all views
			totalViews[data.ArticleID] = 1
			processed := ProcessedEvent{
				ArticleID:    data.ArticleID,
				ArticleTitle: data.ArticleTitle,
				UniqueViews:  uniqueViews[data.ArticleID],
				TotalViews:   totalViews[data.ArticleID],
			}
			// create initial entry for article
			mappedArt[data.ArticleID] = processed
		}

	}

	return mappedArt, nil
}

// GetCountedPosts iterates over a map of articles and retrieves the Articles
// Returns a slice of Article items with updated stats
func GetCountedPosts(data map[string]ProcessedEvent) []ProcessedEvent {
	var articles []ProcessedEvent
	for _, article := range data {
		articles = append(articles, article)
	}
	return articles
}
