package publish

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/stretchr/testify/assert"
)

type mockSNSClient struct {
	snsiface.SNSAPI
}

type mockSNSClientError struct {
	snsiface.SNSAPI
}

func (mockClient *mockSNSClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{
		MessageId: aws.String("test_id"),
	}, nil
}

func (mockClient *mockSNSClientError) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{}, fmt.Errorf("Unable to publish message")
}

func TestSendEventHomePageView(t *testing.T) {
	tag := Tag{Name: "homepage_view"}
	mockSNS := &mockSNSClient{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.Nil(t, err, "Publish successful")
}

func TestSendEventHomePageViewFail(t *testing.T) {
	tag := Tag{Name: "homepage_view"}
	mockSNS := &mockSNSClientError{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.NotNil(t, err, "Publish failed")
}

func TestSendEventPostView(t *testing.T) {
	tag := Tag{Name: "post_view"}
	mockSNS := &mockSNSClient{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.Nil(t, err, "Publish successful")
}

func TestSendEventPostViewFail(t *testing.T) {
	tag := Tag{Name: "post_view"}
	mockSNS := &mockSNSClientError{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.NotNil(t, err, "Publish successful")
}

func TestSendEventContactView(t *testing.T) {
	tag := Tag{Name: "contact_view"}
	mockSNS := &mockSNSClient{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.Nil(t, err, "Publish successful")
}

func TestSendEventContactViewFail(t *testing.T) {
	tag := Tag{Name: "contact_view"}
	mockSNS := &mockSNSClientError{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.NotNil(t, err, "Publish successful")
}

func TestSendEventAboutView(t *testing.T) {
	tag := Tag{Name: "about_view"}
	mockSNS := &mockSNSClient{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.Nil(t, err, "Publish successful")
}

func TestSendEventAboutViewFail(t *testing.T) {
	tag := Tag{Name: "about_view"}
	mockSNS := &mockSNSClientError{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.NotNil(t, err, "Publish successful")
}

func TestSendEventFail(t *testing.T) {
	tag := Tag{Name: "random_view"}
	mockSNS := &mockSNSClientError{}
	err := tag.SendEvent(mockSNS, "testData")
	assert.NotNil(t, err, "Returns an error")
}
