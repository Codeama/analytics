// Package tag tags received events
package tag

import (
	"fmt"
	"regexp"
)

// ViewData represents input data expected
type ViewData struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	ConnectionID string
}

// ProcessedData represents data to be returned
type ProcessedData struct {
	ArticleID    string
	ArticleTitle string
	Event        string
	ConnectionID string
}

func tag(data ViewData, eventTag string) (ProcessedData, error) {
	return ProcessedData{
		ArticleID:    data.ArticleID,
		ArticleTitle: data.ArticleTitle,
		Event:        eventTag,
		ConnectionID: data.ConnectionID,
	}, nil
}

// TranslateData processes and tags events received
func TranslateData(data ViewData) (ProcessedData, error) {
	currentURL := data.CurrentPage
	post, _ := regexp.MatchString("\\/posts\\/", currentURL)
	contact, _ := regexp.MatchString("\\/pages\\/contacts", currentURL)
	about, _ := regexp.MatchString("\\/pages\\/about", currentURL)

	if post {
		result, _ := tag(data, "post_view")
		return result, nil
	}

	if contact {
		result, _ := tag(data, "contact_view")
		return result, nil
	}

	if about {
		result, _ := tag(data, "about_view")
		return result, nil
	}

	err := fmt.Errorf("Unrecognised data format: %s", data)

	return ProcessedData{}, err
}
