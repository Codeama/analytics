package store

import (
	"os"
	"testing"
)

func TestReferrerIsDomain(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		domain   string
		referrer string
		want     bool
	}{
		{
			domain:   "https://example.info/",
			referrer: "https://example.info/posts/random-article",
			want:     true,
		},
		{
			domain:   "https://example.here",
			referrer: "https://example.info/posts/random-article",
			want:     false,
		},
	}

	for _, tc := range testCases {
		os.Setenv("DOMAIN_NAME", tc.domain)
		got := IsDomain(tc.referrer)

		if got != tc.want {
			t.Errorf("IsDomain(%s): want: %v, got: %v", tc.referrer, tc.want, got)
		}
	}
}
