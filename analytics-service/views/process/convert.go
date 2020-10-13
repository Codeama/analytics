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

// ReceivedData represents input data expected
type ReceivedData struct {
	ArticleID    string
	ArticleTitle string
	PreviousPage string
	CurrentPage  string
	ConnectionID string
}

// FilterData filters incoming json data into the right struct
// and returns the corresponding struct as an interface of Event
// for further processing by Sort func
func FilterData(data ReceivedData) Event {
	if data.ArticleID == "" && data.ArticleTitle == "" {
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
func Sort(data Event) (string, string) {
	switch data.(type) {
	case Page:
		currentURL := data.(Page).CurrentPage
		contact, _ := regexp.MatchString("\\/pages\\/contacts", currentURL)
		about, _ := regexp.MatchString("\\/pages\\/about", currentURL)
		home, _ := regexp.MatchString("\\/", currentURL)

		if !contact && !about && !home {
			fmt.Printf("Unrecognised URL from data received: %v \n", currentURL)
			fmt.Printf("Malformed data received: %v", data)
		}

		if about {
			tag, pageData := data.tagEvent("about_view")
			return tag, string(pageData)
		}

		if home {
			tag, pageData := data.tagEvent("homepage_view")
			return tag, string(pageData)
		}

		if contact {
			tag, pageData := data.tagEvent("contact_view")
			return tag, string(pageData)
		}

	case Article:
		currentURL := data.(Article).CurrentPage
		post, _ := regexp.MatchString("\\/posts\\/", currentURL)
		if post {
			tag, pageData := data.tagEvent("post_view")
			return tag, string(pageData)
		}
	default:
		fmt.Printf("Cannot process unknow data type %v", data)
	}
	return "", ""
}
