package process

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var mockCountResult = map[string]ProcessedEvent{
	"testArticleID": ProcessedEvent{
		ArticleID:    "testArticleID",
		ArticleTitle: "Test Title",
		UniqueViews:  2,
		TotalViews:   3,
	},
	"testArticleID2": ProcessedEvent{
		ArticleID:    "testArticleID2",
		ArticleTitle: "Test Title",
		UniqueViews:  1,
		TotalViews:   1,
	},
	"testArticleID3": ProcessedEvent{
		ArticleID:    "testArticleID3",
		ArticleTitle: "Test Title",
		UniqueViews:  3,
		TotalViews:   3,
	},
	"testArticleID4": ProcessedEvent{
		ArticleID:    "testArticleID4",
		ArticleTitle: "Test Title",
		UniqueViews:  1,
		TotalViews:   1,
	},
	"testArticleID5": ProcessedEvent{
		ArticleID:    "testArticleID5",
		ArticleTitle: "Test Title",
		UniqueViews:  1,
		TotalViews:   2,
	},
}

// 1. It should count all views for each article and return an array of the articles and their stats
// 2. It should find unique views for each article and count them
func TestCountViews(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("../testdata/article-events.json")
	if err != nil {
		t.Errorf("could not read test data")
	}
	var inputEvent events.SQSEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal data. details: %v", err)
	}
	mapData, _ := CountViews(inputEvent)
	assert.Equal(t, mockCountResult, mapData)
}
