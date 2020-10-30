package process

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// Data used to mock test result for TestCountView
// Also used as data input for TestGetPost
var mockViewCount = map[string]ProcessedEvent{
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

// Data used as TestGetPosts expected result
var mockPostResult = []ProcessedEvent{
	{
		ArticleID:    "testArticleID",
		ArticleTitle: "Test Title",
		UniqueViews:  2,
		TotalViews:   3,
	},
	{
		ArticleID:    "testArticleID2",
		ArticleTitle: "Test Title",
		UniqueViews:  1,
		TotalViews:   1,
	},
	{
		ArticleID:    "testArticleID3",
		ArticleTitle: "Test Title",
		UniqueViews:  3,
		TotalViews:   3,
	},
	{
		ArticleID:    "testArticleID4",
		ArticleTitle: "Test Title",
		UniqueViews:  1,
		TotalViews:   1,
	},
	{
		ArticleID:    "testArticleID5",
		ArticleTitle: "Test Title",
		UniqueViews:  1,
		TotalViews:   2,
	},
}

// 1. It should count all views (total and unique) for each article and return an array/slice of the articles and their stats
func TestCountViews(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("../testdata/article-events.json")
	if err != nil {
		t.Errorf("could not read test data")
	}
	var inputEvent events.SQSEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal data. details: %v", err)
	}
	actual, _ := countViews(inputEvent)
	assert.Equal(t, mockViewCount, actual)
}

func TestGetPosts(t *testing.T) {
	var expected ProcessedEvent
	posts := GetPosts(mockViewCount)
	assert.Equal(t, len(mockPostResult), len(posts), "It should return an array of same length as the map input data")
	// no guarantee of map item order so iterating over both actual and expected
	for _, actual := range posts {
		for _, article := range mockPostResult {
			if actual == article {
				expected = article
				break
			}
		}
		assert.Equal(t, actual, expected, "Map input values should be the same as the returned slice values")
	}
}
