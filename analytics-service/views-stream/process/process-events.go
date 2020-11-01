// Package process tags received events
package process

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Event contains functions for each received event
type Event interface {
	tagEvent(tag string) (string, string)
}

// Page represent structure all incoming data should have
type Page struct {
	PreviousPage string
	CurrentPage  string
	ConnectionID string
}

// Article represents post view data
type Article struct {
	ArticleID    string
	ArticleTitle string
	Page
}

// IncomingData represents data event received
type IncomingData struct {
	ArticleID    string `json:"articleId,omitempty"`
	ArticleTitle string `json:"articleTitle,omitempty"`
	PreviousPage string `json:"previousPage"`
	CurrentPage  string `json:"currentPage"`
}

// AnalyticsData represents input data expected
type AnalyticsData struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	ConnectionID string
}

// ValidateData validates and stores incoming data in an struct
func ValidateData(data IncomingData, id string) (AnalyticsData, error) {
	// check connectionId isn't empty
	if id == "" {
		return AnalyticsData{}, fmt.Errorf("No ConnectionId from request")
	}

	// check there are at least members of Page struct
	if data.CurrentPage == "" || data.PreviousPage == "" {
		return AnalyticsData{}, fmt.Errorf("Event does not contain required page data")
	}

	// return received data in a digestable struct
	return AnalyticsData{
		data.ArticleID,
		data.ArticleTitle,
		data.PreviousPage,
		data.CurrentPage,
		id,
	}, nil

}

// FilterData filters incoming json data into the right struct
// and returns the corresponding struct as an interface of Event
// for further processing by Sort func
func FilterData(data AnalyticsData) Event {
	if data.ArticleID == "" || data.ArticleTitle == "" {
		var page Page
		page.ConnectionID = data.ConnectionID
		page.PreviousPage = data.PreviousPage
		page.CurrentPage = data.CurrentPage
		return page
	}

	var article Article
	article.ArticleID = data.ArticleID
	article.ArticleTitle = data.ArticleTitle
	article.PreviousPage = data.PreviousPage
	article.CurrentPage = data.CurrentPage
	article.ConnectionID = data.ConnectionID
	return article
}

// Tags received events that match type Page
func (data Page) tagEvent(eventTag string) (string, string) {
	page := struct {
		ConnectionID string
		CurrentPage  string
		PreviousPage string
		EventType    string
	}{
		data.ConnectionID,
		data.CurrentPage,
		data.PreviousPage,
		eventTag,
	}
	result, _ := json.Marshal(page)
	return eventTag, string(result)
}

// Tags received events that match type Article
func (data Article) tagEvent(eventTag string) (string, string) {
	post := struct {
		ArticleID    string
		ArticleTitle string
		ConnectionID string
		CurrentPage  string
		PreviousPage string
		EventType    string
	}{
		data.ArticleID,
		data.ArticleTitle,
		data.Page.ConnectionID,
		data.Page.CurrentPage,
		data.Page.PreviousPage,
		eventTag,
	}
	result, _ := json.Marshal(post)
	return eventTag, string(result)
}

// Sort processes and tags received events for publishing to SNS
// Identifies page url, tags it, marshals it and
// returns the tag and json data (eventTag, data)
func Sort(data Event) (string, string, error) {
	switch data.(type) {
	case Page:
		currentURL := data.(Page).CurrentPage
		contact := currentURL == "/pages/contacts"
		about := currentURL == "/pages/about"
		home := currentURL == "/"

		if about {
			tag, pageData := data.tagEvent("about_view")
			return tag, string(pageData), nil
		}

		if home {
			tag, pageData := data.tagEvent("homepage_view")
			return tag, string(pageData), nil
		}

		if contact {
			tag, pageData := data.tagEvent("contact_view")
			return tag, string(pageData), nil
		}

		return "", "", fmt.Errorf("Unrecognised URL %v.\n Data received: %v", currentURL, data)

	case Article:
		currentURL := data.(Article).CurrentPage
		post, _ := regexp.MatchString("\\/posts\\/", currentURL)
		if post {
			tag, pageData := data.tagEvent("post_view")
			return tag, string(pageData), nil
		}
		return "", "", fmt.Errorf("Unrecognised URL %v.\n Data received: %v", currentURL, data)
	}
	return "", "", fmt.Errorf("Unknown event %v", data)
}
