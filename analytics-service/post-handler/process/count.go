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
	EventType    string
	ConnectionID string
}

// ProcessedEvent is the result returned after counting
type ProcessedEvent struct {
	ArticleID    string
	ArticleTitle string
	UniqueViews  int
	TotalViews   int // includes a count of non-unique views
}

// Checks if event is a unique view of post/article
func isUnique(previousPage string, currentPage string) bool {
	// UNIQUE: if currentPage is not a / and previousPage is not null
	if previousPage != "null" {
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
				UniqueViews:  uniqueViews[data.ArticleID],
				TotalViews:   totalViews[data.ArticleID],
			}
			// update article values
			mappedArt[data.ArticleID] = processed
		} else {
			unique := isUnique(data.PreviousPage, data.CurrentPage)
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
func GetCountedPosts(data map[string]ProcessedEvent) []ProcessedEvent {
	var articles []ProcessedEvent
	for _, article := range data {
		articles = append(articles, article)
	}
	return articles
}
