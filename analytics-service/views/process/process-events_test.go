package process

import "testing"

var jsonArticleMock = `{
	"articleId":"123testId",
	"articleTitle": "Unit Testing Go Functions",
	"previousPage": "/",
	"currentPage": "/posts/unit-testing-go-functions"
	}`

var mockArticleData = ReceivedData{
	ArticleID:    "123testId",
	ArticleTitle: "Unit Testing Go Functions",
	PreviousPage: "/",
	CurrentPage:  "/posts/unit-testing-go-functions",
}

var mockPageData = ReceivedData{
	PreviousPage: "null",
	CurrentPage:  "/",
}

func TestFilterDataArticle(t *testing.T) {
	// test it returns an Article struct
	result := FilterData(mockArticleData)
	if result != result.(Event) {
		t.Errorf("Expected result to be type %T, but received %T", mockArticleData, result)
	}

}

func TestFilterDataPage(t *testing.T) {
	// test it returns a Page struct
	result := FilterData(mockPageData)
	if result != result.(Event) {
		t.Errorf("Expected result to be type %T, but received %T", mockPageData, result)
	}

}

func TestSort(t *testing.T) {}
