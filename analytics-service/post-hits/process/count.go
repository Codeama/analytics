package process

import (
	"encoding/json"
	"fmt"
	"os"

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
	Referrer     string
}

// ProcessedEvent is the result returned after counting
type ProcessedEvent struct {
	ArticleID    string
	ArticleTitle string
	UniqueViews  int
	TotalViews   int // sum total of all views unique or not
}

// Checks if event is a unique view of post/article
func isUnique(previousPage string, referrer string) bool {
	// UNIQUE: if previousPage is not null OR if previousPage is null and referrer is not current domain
	if previousPage != "null" || previousPage == "null" && referrer != os.Getenv("DOMAIN_NAME") {
		return true
	}

	return false
}

// CountViews totals the number of views for each article
// It processes each article event by implementing a set using Go map
// and totals them, each one being a view (1)
// Employs a map reduce paradigm for parallel data processing
func CountViews(sqsEvent events.SQSEvent) (map[string]ProcessedEvent, error) {
	var data incomingEvent
	var totalViews = make(map[string]int)
	var uniqueViews = make(map[string]int)
	var mappedArt = make(map[string]ProcessedEvent)
	for _, message := range sqsEvent.Records {
		// serialise to Go struct
		if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
			// fmt.Println("Could not deserialise data: ", err)
			return nil, fmt.Errorf("Could not deserialise data: %v", err)
		}
		// checks current article has a view value
		_, hasViews := totalViews[data.ArticleID]
		// checks article is already in map (it should be; see else statement below that runs at least once for all inicoming data)
		// then updates it (deletes and replaces with article item
		// with the latest view count)
		_, exists := mappedArt[data.ArticleID]
		if hasViews && exists {
			delete(mappedArt, data.ArticleID)

			// process unique views
			unique := isUnique(data.PreviousPage, data.Referrer)
			if unique {
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
			unique := isUnique(data.PreviousPage, data.Referrer)
			// process unique views
			if unique {
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
