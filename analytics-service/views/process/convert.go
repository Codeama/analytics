// Package process tags received events
package process

import (
	"encoding/json"
	"fmt"
	"regexp"
)

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

// Received data can be Page or Article

// OutgoingData represents data to be returned
type OutgoingData struct {
	ArticleID    string
	ArticleTitle string
	Event        string
	ConnectionID string
}

// FilterData filters incoming json data into the right struct
// for further processing by Sort func
func FilterData(data ReceivedData) interface{} {
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

// Sort processes and tags received events for publishing to SNS
// returns (eventTag, data)
func Sort(data interface{}) (string, string) {
	switch data.(type) {
	case Page:
		currentURL := data.(Page).CurrentPage
		contact, _ := regexp.MatchString("\\/pages\\/contacts", currentURL)
		about, _ := regexp.MatchString("\\/pages\\/about", currentURL)
		home, _ := regexp.MatchString("\\/", currentURL)

		// TODO extract to a new function and/or use switch
		if !contact && !about && !home {
			fmt.Printf("Unrecognised URL from data received: %v \n", currentURL)
			fmt.Printf("Malformed data received: %v", data)
		}

		if home {
			tag := "homepage_view"
			homePage := struct {
				ConnectionID string
				CurrentPage  string
				PreviousPage string
				EventType    string
			}{
				data.(Page).ConnectionID,
				data.(Page).CurrentPage,
				data.(Page).PreviousPage,
				tag,
			}
			result, _ := json.Marshal(homePage)
			return tag, string(result)
		}
		if contact {
			tag := "contact_view"
			contactMe := struct {
				ConnectionID string
				CurrentPage  string
				PreviousPage string
				EventType    string
			}{
				data.(Page).ConnectionID,
				data.(Page).CurrentPage,
				data.(Page).PreviousPage,
				tag,
			}
			result, _ := json.Marshal(contactMe)
			return tag, string(result)
		}

		if about {
			tag := "about_view"
			aboutMe := struct {
				ConnectionID string
				CurrentPage  string
				PreviousPage string
				EventType    string
			}{
				data.(Page).ConnectionID,
				data.(Page).CurrentPage,
				data.(Page).PreviousPage,
				tag,
			}
			result, _ := json.Marshal(aboutMe)
			return tag, string(result)
		}
	case Article:
		currentURL := data.(Article).CurrentPage
		post, _ := regexp.MatchString("\\/posts\\/", currentURL)
		if post {
			tag := "post_view"
			post := struct {
				ArticleID    string
				ArticleTitle string
				ConnectionID string
				CurrentPage  string
				PreviousPage string
				EventType    string
			}{
				data.(Article).ArticleID,
				data.(Article).ArticleTitle,
				data.(Article).Page.ConnectionID,
				data.(Article).Page.CurrentPage,
				data.(Article).Page.PreviousPage,
				tag,
			}
			result, _ := json.Marshal(post)
			return tag, string(result)
		}
	default:
		fmt.Printf("Cannot process unknow data type %v", data)
		// return "", ""
	}
	return "", ""
}
