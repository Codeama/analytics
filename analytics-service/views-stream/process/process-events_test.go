package process

import (
	"testing"
)

func TestValidateData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input1      IncomingData
		input2      string
		want        AnalyticsData
		errReturned bool
	}{
		{
			input1:      IncomingData{ArticleID: "TestID", ArticleTitle: "TestArticle", PreviousPage: "/mypage", CurrentPage: "/posts/hello-world", Refreshed: false},
			input2:      "testId",
			want:        AnalyticsData{ArticleID: "TestID", ArticleTitle: "TestArticle", PreviousPage: "/mypage", CurrentPage: "/posts/hello-world", ConnectionID: "testId", Refreshed: false},
			errReturned: false,
		},
		{
			input1:      IncomingData{ArticleID: "TestID", ArticleTitle: "TestArticle", PreviousPage: "/", CurrentPage: "/posts/hello-world", Refreshed: false},
			input2:      "",
			want:        AnalyticsData{ArticleID: "TestID", ArticleTitle: "TestArticle", PreviousPage: "/", CurrentPage: "/page/unnecessary", ConnectionID: "testID", Refreshed: false},
			errReturned: true,
		},
		{
			input1:      IncomingData{ArticleID: "", ArticleTitle: "", PreviousPage: "/", CurrentPage: "/posts/hello-world", Refreshed: false},
			input2:      "",
			want:        AnalyticsData{ArticleID: "TestID", ArticleTitle: "TestArticle", PreviousPage: "/", CurrentPage: "/page/unnecessary", ConnectionID: "testID", Refreshed: false},
			errReturned: true,
		},
	}

	for _, tc := range testCases {
		got, err := ValidateData(tc.input1, tc.input2)
		errExpected := (err != nil)

		if tc.errReturned != errExpected {
			t.Fatalf("ValidateData(%v, %s): unexpected error status %v", tc.input1, tc.input2, errExpected)
		}
		if !errExpected && got != tc.want {
			t.Errorf("ValidateData(%v, %s) expected: %v, got: %v", tc.input1, tc.input2, tc.want, got)
		}
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input AnalyticsData
		want  Event
	}{
		{
			name:  "Article",
			input: AnalyticsData{ArticleID: "123testId", ArticleTitle: "Unit Testing Go Functions", PreviousPage: "/", CurrentPage: "/posts/unit-testing-go-functions", Refreshed: false},
			want:  Article{"123testId", "Unit Testing Go Functions", Page{PreviousPage: "/", CurrentPage: "/posts/unit-testing-go-functions", Refreshed: false}},
		},
		{
			name:  "Page",
			input: AnalyticsData{ArticleID: "", ArticleTitle: "", PreviousPage: "/", CurrentPage: "/posts/unit-testing-go-functions", Refreshed: false},
			want:  Page{PreviousPage: "/", CurrentPage: "/posts/unit-testing-go-functions", Refreshed: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FilterData(tc.input)

			if got != tc.want {
				t.Errorf("FilterData(%v) expected: %v, got: %v", tc.input, tc.want, got)
			}
		})
	}
}

func TestSort(t *testing.T) {
	testCases := []struct {
		name        string
		input       Event
		output1     string
		output2     string
		errReturned bool
	}{
		{
			name:        "Homepage",
			input:       Page{PreviousPage: "null", CurrentPage: "/", Refreshed: false},
			output1:     "homepage_view",
			output2:     string(`{"ConnectionID":"","CurrentPage":"/","PreviousPage":"null","Refreshed":false,"Referrer":"","EventType":"homepage_view"}`),
			errReturned: false,
		},
		{
			name:        "AboutPage",
			input:       Page{PreviousPage: "/", CurrentPage: "/pages/about", Refreshed: false},
			output1:     "about_view",
			output2:     string(`{"ConnectionID":"","CurrentPage":"/pages/about","PreviousPage":"/","Refreshed":false,"Referrer":"","EventType":"about_view"}`),
			errReturned: false,
		},
		{
			name:        "ContactsPage",
			input:       Page{PreviousPage: "/", CurrentPage: "/pages/contacts", Refreshed: false},
			output1:     "contact_view",
			output2:     string(`{"ConnectionID":"","CurrentPage":"/pages/contacts","PreviousPage":"/","Refreshed":false,"Referrer":"","EventType":"contact_view"}`),
			errReturned: false,
		},
		{
			name:        "ArticlePage",
			input:       Article{"123testId", "Unit Testing Go Functions", Page{PreviousPage: "/", CurrentPage: "/posts/unit-testing-go-functions", Refreshed: false}},
			output1:     "post_view",
			output2:     string(`{"ArticleID":"123testId","ArticleTitle":"Unit Testing Go Functions","ConnectionID":"","CurrentPage":"/posts/unit-testing-go-functions","PreviousPage":"/","Refreshed":false,"Referrer":"","EventType":"post_view"}`),
			errReturned: false,
		},
		{
			name:        "Error",
			input:       Page{PreviousPage: "/", CurrentPage: "/me/about", Refreshed: false},
			output1:     "bogusText",
			output2:     "BOGUS",
			errReturned: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1, got2, err := Sort(tc.input)

			errExpected := (err != nil)

			if tc.errReturned != errExpected {
				t.Fatalf("Sort(%v) unexpected error status: want %v got %v", tc.input, tc.errReturned, err)
			}
			if !errExpected && got1 != tc.output1 {
				t.Errorf("Sort(%v) expected %v, got %v", tc.input, tc.output1, got1)
			}
			if !errExpected && got2 != tc.output2 {
				t.Errorf("Sort(%v) expected %v, got %v", tc.input, tc.output2, got2)
			}
		})
	}
}

func TestUnknownEventType(t *testing.T) {
	tc := struct {
		input       mockEventType
		want        string
		errReturned bool
	}{
		input:       mockEventType{"hello-world-page"},
		want:        "Doesn't matter when it returns an error",
		errReturned: true,
	}

	_, _, err := Sort(tc.input)

	errExpected := (err != nil)

	if errExpected != tc.errReturned {
		t.Fatalf("Sort(%v) unexpected error status: expected %v got %v", tc.input, tc.errReturned, err)
	}
}

// Implements Event interface to test for unknown event
type mockEventType struct {
	testPage string
}

func (mockedData mockEventType) tagEvent(tag string) (string, string) {
	return "test_tag", "testData"
}
