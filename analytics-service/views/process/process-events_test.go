package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockJSONArticle = `{
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

var mockHomePageData = AnalyticsData{
	PreviousPage: "null",
	CurrentPage:  "/",
}

var mockHomePage = Page{
	PreviousPage: "null",
	CurrentPage:  "/",
}

var mockAboutMePage = Page{
	PreviousPage: "/",
	CurrentPage:  "/pages/about",
}

var mockContactMePage = Page{
	PreviousPage: "/",
	CurrentPage:  "/pages/contacts",
}

var mockPageDataError = Page{
	PreviousPage: "/",
	CurrentPage:  "/me/about",
}

var mockArticlePage = Article{
	ArticleID:    "123testId",
	ArticleTitle: "Unit Testing Go Functions",
	Page: Page{
		PreviousPage: "/",
		CurrentPage:  "/posts/unit-testing-go-functions",
	},
}

var mockArticlePageError = Article{
	ArticleID:    "123testId",
	ArticleTitle: "Random Title",
	Page: Page{
		PreviousPage: "/",
		CurrentPage:  "/posts/",
	},
}

func TestFilterDataArticle(t *testing.T) {
	result := FilterData(mockArticleData)
	assert.Equal(t, result.(Event), result, "Article should be of type Event")

}

func TestFilterDataPage(t *testing.T) {
	result := FilterData(mockHomePageData)
	assert.Equal(t, result.(Event), result, "Page should be of type Event")

}

func TestSortHome(t *testing.T) {
	var mockJSONHomePage = `{"ConnectionID":"","CurrentPage":"/","PreviousPage":"null","EventType":"homepage_view"}`
	tag, data, _ := Sort(mockHomePage)
	assert.Equal(t, "homepage_view", tag, "Event tag should be 'homepage_view'")
	assert.Equal(t, string(mockJSONHomePage), data, "It should return a JSON string")
}

func TestSortAbout(t *testing.T) {
	var mockJSONAboutMePage = `{"ConnectionID":"","CurrentPage":"/pages/about","PreviousPage":"/","EventType":"about_view"}`
	tag, data, _ := Sort(mockAboutMePage)
	assert.Equal(t, "about_view", tag, "Event tag should be 'about_view'")
	assert.Equal(t, string(mockJSONAboutMePage), data, "It should return a JSON string")
}

func TestSortContact(t *testing.T) {
	var mockJSONContactMePage = `{"ConnectionID":"","CurrentPage":"/pages/contacts","PreviousPage":"/","EventType":"contact_view"}`
	tag, data, _ := Sort(mockContactMePage)
	assert.Equal(t, "contact_view", tag, "Event tag should be 'contact_view'")
	assert.Equal(t, string(mockJSONContactMePage), data, "It should return a JSON string")
}

func TestSortArticle(t *testing.T) {
	var mockJSONArticlePage = `{"ArticleID":"123testId","ArticleTitle":"Unit Testing Go Functions","ConnectionID":"","CurrentPage":"/posts/unit-testing-go-functions","PreviousPage":"/","EventType":"post_view"}`
	tag, data, _ := Sort(mockArticlePage)
	assert.Equal(t, "post_view", tag, "Event tag should be 'post_view'")
	assert.Equal(t, string(mockJSONArticlePage), data, "It should return a JSON string")
}

func TestSortPageError(t *testing.T) {
	tag, data, err := Sort(mockPageDataError)
	assert.Equal(t, "", tag, "There should be no tag")
	assert.Equal(t, "", data, "There should be no data")
	assert.NotNil(t, err, "There should be an error")
}

func TestSortArticleError(t *testing.T) {
	tag, data, err := Sort(mockPageDataError)
	assert.Equal(t, "", tag, "There should be no tag")
	assert.Equal(t, "", data, "There should be no data")
	assert.NotNil(t, err, "There should be an error")
}
