package rss

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockFeedFetcher is a mock implementation of FeedFetcher for testing.
type MockFeedFetcher struct {
	Items []RSSInfo
	Err   error
}

// Fetch returns the mock feed items or an error.
func (m *MockFeedFetcher) Fetch(feedURL string) ([]RSSInfo, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Items, nil
}

func TestFetchRSS(t *testing.T) {
	mockFetcher := &MockFeedFetcher{
		Items: []RSSInfo{
			{Title: "Title 1", Link: "http://example.com/1"},
			{Title: "Title 2", Link: "http://example.com/2"},
		},
	}

	result, err := fetchRSS(mockFetcher, "http://example.com/feed")
	assert.NoError(t, err)
	assert.Equal(t, mockFetcher.Items, result)
}

func TestFetchRSS_Error(t *testing.T) {
	mockFetcher := &MockFeedFetcher{
		Err: errors.New("fetch error"),
	}

	_, err := fetchRSS(mockFetcher, "http://example.com/feed")
	assert.Error(t, err)
	assert.Equal(t, "failed to fetch RSS feed from URL http://example.com/feed: fetch error", err.Error())
}
