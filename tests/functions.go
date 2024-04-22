package tests

import (
	"net/url"
	"testing"
)

func MustParseURL(t *testing.T, s string) *url.URL {
	result, err := url.Parse(s)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}
	return result
}
