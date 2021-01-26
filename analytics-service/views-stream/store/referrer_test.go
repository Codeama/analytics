package store

import (
	"fmt"
	"os"
	"testing"
)

func TestReferrerIsDomain(t *testing.T) {
	t.Parallel()
	os.Setenv("DOMAIN_NAME", "https://example.info/")
	fmt.Println(os.Getenv("DOMAIN_NAME"))
	fakeReferrer := "https://example.info/posts/random-article"
	expected := true
	actual := IsDomain(fakeReferrer)

	if actual != expected {
		t.Errorf("IsDomain(%s): want: %v, got: %v", fakeReferrer, expected, actual)
	}

}

func TestReferrerIsNotDomain(t *testing.T) {
	os.Setenv("DOMAIN_NAME", "https://example.here")
	fmt.Println(os.Getenv("DOMAIN_NAME"))
	fakeReferrer := "https://example.info/posts/random-article"
	expected := false
	actual := IsDomain(fakeReferrer)

	if actual != expected {
		t.Errorf("IsDomain(%s): want: %v, got: %v", fakeReferrer, expected, actual)
	}
}
