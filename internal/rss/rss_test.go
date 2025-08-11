package rss

import (
	"errors"
	"testing"

	"github.com/mmcdole/gofeed"
)

// MockFeedFetcher is a mock implementation of FeedFetcher for testing.
type MockFeedFetcher struct {
	Feed *gofeed.Feed
	Err  error
}

// Fetch returns the mock feed or an error.
func (m *MockFeedFetcher) Fetch(feedUrl string) (*gofeed.Feed, error) {
	return m.Feed, m.Err
}

func TestFetchRss(t *testing.T) {
	mockFeed := &gofeed.Feed{
		Items: []*gofeed.Item{
			{Title: "Item 1"},
			{Title: "Item 2"},
			{Title: "Item 3"},
		},
	}
	mockFetcher := &MockFeedFetcher{Feed: mockFeed}

	result, err := fetchRss(mockFetcher, "https://example.com/feed")
	if err != nil {
		t.Fatalf("FetchRss returned an error: %v", err)
	}

	expected := `[{"title":"Item 1"},{"title":"Item 2"},{"title":"Item 3"}]`
	if result != expected {
		t.Errorf("FetchRss() = %v; want %v", result, expected)
	}
}

func TestFetchRss_FetchError(t *testing.T) {
	mockFetcher := &MockFeedFetcher{Err: errors.New("failed to fetch")}

	_, err := fetchRss(mockFetcher, "https://example.com/feed")
	if err == nil {
		t.Fatal("FetchRss() expected an error but got nil")
	}
}
