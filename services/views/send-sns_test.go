package views

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// Tests in this file follow the testing advice given on AWS docs/repos
// which helps to mock their services so I can test my own code
type mockSNSClient struct {
	snsiface.SNSAPI
}

func (mockClient *mockSNSClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{
		MessageId: aws.String("test_id"),
	}, nil
}

func TestSendEventTags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		tag         Tag
		event       string
		errReturned bool
	}{
		{
			name:        "HomePageView",
			tag:         Tag{"homepage_view"},
			event:       "testData",
			errReturned: false,
		},
		{
			name:        "PostView",
			tag:         Tag{"post_view"},
			event:       "from-client-data",
			errReturned: false,
		},
		{
			name:        "UnknownTag",
			tag:         Tag{"just_mine"},
			event:       "asdf;lkj",
			errReturned: true,
		},
		{
			name:        "ContactView",
			tag:         Tag{"contact_view"},
			event:       "any data",
			errReturned: false,
		},
		{
			name:        "AboutView",
			tag:         Tag{"about_view"},
			event:       "asdf;lkj",
			errReturned: false,
		},
		{
			name:        "Random",
			tag:         Tag{"randoms"},
			event:       "any data",
			errReturned: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.tag.SendEvent(&mockSNSClient{}, tc.event)

			errExpected := (err != nil)

			if tc.errReturned != errExpected {
				t.Fatalf("Unexpected error status: expected no error, got %q", err)
			}
		})
	}
}
