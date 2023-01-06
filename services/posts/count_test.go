package posts_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/codeama/analytics/services/posts"
)

// Data used to mock test result for TestCountView
// Also used as data input for TestGetPost
var countResult = map[string]posts.ProcessedEvent{
	"testArticleID": posts.ProcessedEvent{
		ArticleID:    "testArticleID",
		ArticleTitle: "Test Title 1",
		UniqueViews:  1,
		TotalViews:   4,
	},
	"testArticleID2": posts.ProcessedEvent{
		ArticleID:    "testArticleID2",
		ArticleTitle: "Test Title 2",
		UniqueViews:  0,
		TotalViews:   1,
	},
	"testArticleID3": posts.ProcessedEvent{
		ArticleID:    "testArticleID3",
		ArticleTitle: "Test Title 3",
		UniqueViews:  2,
		TotalViews:   3,
	},
	"testArticleID4": posts.ProcessedEvent{
		ArticleID:    "testArticleID4",
		ArticleTitle: "Test Title 4",
		UniqueViews:  1,
		TotalViews:   1,
	},
	"testArticleID5": posts.ProcessedEvent{
		ArticleID:    "testArticleID5",
		ArticleTitle: "Test Title 5",
		UniqueViews:  1,
		TotalViews:   2,
	},
}

// Data used as TestGetPosts expected result
var getPostResult = []posts.ProcessedEvent{
	{
		ArticleID:    "testArticleID",
		ArticleTitle: "Test Title 1",
		UniqueViews:  2,
		TotalViews:   4,
	},
	{
		ArticleID:    "testArticleID2",
		ArticleTitle: "Test Title 2",
		UniqueViews:  0,
		TotalViews:   1,
	},
	{
		ArticleID:    "testArticleID3",
		ArticleTitle: "Test Title 3",
		UniqueViews:  2,
		TotalViews:   3,
	},
	{
		ArticleID:    "testArticleID4",
		ArticleTitle: "Test Title 4",
		UniqueViews:  1,
		TotalViews:   1,
	},
	{
		ArticleID:    "testArticleID5",
		ArticleTitle: "Test Title 5",
		UniqueViews:  1,
		TotalViews:   2,
	},
}

// 1. It should count all views (total and unique) for each article and return an array/slice of the articles and their stats
func TestCountViews(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/article-events.json")
	if err != nil {
		t.Errorf("could not read test data")
	}
	var inputEvent events.SQSEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal data. details: %v", err)
	}
	got, err := posts.CountViews(inputEvent)

	if err != nil {
		t.Fatalf("Expected no error, got %q", err)
	}

	if !cmp.Equal(got, countResult) {
		t.Errorf(cmp.Diff(countResult, got))
	}
}

func TestGetPosts(t *testing.T) {
	var expectedProcessed posts.ProcessedEvent
	var actualProcessed posts.ProcessedEvent
	posts := posts.GetCountedPosts(countResult)

	if len(getPostResult) != len(posts) {
		t.Errorf("GetCountedPosts(%v) want total items %v, got %v", countResult, len(getPostResult), len(posts))
	}
	// no guarantee of map item order so iterating over both actual and expected
	// and assigning values in the same map order for both expected and actual
	for _, actual := range posts {
		for _, article := range getPostResult {
			if actual == article {
				actualProcessed = actual
				expectedProcessed = article
				break
			}
		}

		if !cmp.Equal(expectedProcessed, actualProcessed) {
			t.Errorf(cmp.Diff(expectedProcessed, actualProcessed))
		}
	}
}
