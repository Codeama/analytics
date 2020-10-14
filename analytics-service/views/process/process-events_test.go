package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonArticleMock = `{
	"articleId":"123testId",
	"articleTitle": "Unit Testing Go Functions",
	"previousPage": "/",
	"currentPage": "/posts/unit-testing-go-functions"
	}`

var mockArticleData = AnalyticsData{
	ArticleID:    "123testId",
	ArticleTitle: "Unit Testing Go Functions",
	PreviousPage: "/",
	CurrentPage:  "/posts/unit-testing-go-functions",
}

var mockPageData = AnalyticsData{
	PreviousPage: "null",
	CurrentPage:  "/",
}

func TestFilterDataArticle(t *testing.T) {
	result := FilterData(mockArticleData)
	assert.Equal(t, result, result.(Event), "Article should be of type Event")

}

func TestFilterDataPage(t *testing.T) {
	result := FilterData(mockPageData)
	assert.Equal(t, result, result.(Event), "Page should be of type Event")

}

func TestSort(t *testing.T) {

}
