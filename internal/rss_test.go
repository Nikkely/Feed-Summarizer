package rss

import (
	"reflect"
	"testing"
)

func TestFetchRss(t *testing.T) {
	feedUrl := "https://example.com/feed"
	expected := []string{"Item 1", "Item 2", "Item 3"}

	result, err := FetchRss(feedUrl)
	if err != nil {
		t.Fatalf("fetchRss returned an error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("fetchRss(%q) = %v; want %v", feedUrl, result, expected)
	}
}
