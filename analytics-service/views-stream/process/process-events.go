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

// Page represent structure of incoming page view data
type Page struct {
	PreviousPage string
	CurrentPage  string
	ConnectionID string
	Refreshed    bool
	Referrer     string
}

// Article represents post view data
type Article struct {
	ArticleID    string
	ArticleTitle string
	Page
}

// IncomingData represents page view event received
type IncomingData struct {
	ArticleID    string `json:"articleId,omitempty"`
	ArticleTitle string `json:"articleTitle,omitempty"`
	PreviousPage string `json:"previousPage"`
	CurrentPage  string `json:"currentPage"`
	Refreshed    bool   `json:"refreshed"`
	Referrer     string `json:"referrer"`
}

// AnalyticsData is transformed data from IncomingData
type AnalyticsData struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	ConnectionID string
	Refreshed    bool
	Referrer     string
}

// ValidateData validates and stores incoming data and a connection ID
func ValidateData(data IncomingData, id string) (AnalyticsData, error) {
	// check connectionId isn't empty
	if id == "" {
		return AnalyticsData{}, fmt.Errorf("no ConnectionId from request")
	}

	// check there are at least members of Page struct
	if data.CurrentPage == "" || data.PreviousPage == "" {
		return AnalyticsData{}, fmt.Errorf("Event does not contain required page data")
	}

	// map and return received data to a struct
	return AnalyticsData{
		data.ArticleID,
		data.ArticleTitle,
		data.PreviousPage,
		data.CurrentPage,
		id,
		data.Refreshed,
		data.Referrer,
	}, nil

}

// FilterData takes AnalyticsData struct, filters and maps the members into
// an Article or Page struct and returns an Event interface
func FilterData(data AnalyticsData) Event {
	if data.ArticleID == "" || data.ArticleTitle == "" {
		var page Page
		page.ConnectionID = data.ConnectionID
		page.PreviousPage = data.PreviousPage
		page.CurrentPage = data.CurrentPage
		page.Refreshed = data.Refreshed
		page.Referrer = data.Referrer
		return page
	}

	var article Article
	article.ArticleID = data.ArticleID
	article.ArticleTitle = data.ArticleTitle
	article.PreviousPage = data.PreviousPage
	article.CurrentPage = data.CurrentPage
	article.ConnectionID = data.ConnectionID
	article.Refreshed = data.Refreshed
	article.Referrer = data.Referrer
	return article
}

// Tags received events that qualify as Page
func (data Page) tagEvent(eventTag string) (string, string) {
	page := struct {
		ConnectionID string
		CurrentPage  string
		PreviousPage string
		Refreshed    bool
		Referrer     string
		EventType    string
	}{
		data.ConnectionID,
		data.CurrentPage,
		data.PreviousPage,
		data.Refreshed,
		data.Referrer,
		eventTag,
	}
	result, _ := json.Marshal(page)
	return eventTag, string(result)
}

// Tags received events that qualify as Article
func (data Article) tagEvent(eventTag string) (string, string) {
	post := struct {
		ArticleID    string
		ArticleTitle string
		ConnectionID string
		CurrentPage  string
		PreviousPage string
		Refreshed    bool
		Referrer     string
		EventType    string
	}{
		data.ArticleID,
		data.ArticleTitle,
		data.Page.ConnectionID,
		data.Page.CurrentPage,
		data.Page.PreviousPage,
		data.Page.Refreshed,
		data.Page.Referrer,
		eventTag,
	}
	result, _ := json.Marshal(post)
	return eventTag, string(result)
}

// Sort processes and tags received events for publishing to SNS
// Identifies page url, tags it, marshals it and
// returns the tag and json formatted data (eventTag, data)
func Sort(data Event) (string, string, error) {
	switch data.(type) {
	case Page:
		currentURL := data.(Page).CurrentPage
		contact, _ := regexp.MatchString("/pages/contacts", currentURL)
		about, _ := regexp.MatchString("/pages/about", currentURL)
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
