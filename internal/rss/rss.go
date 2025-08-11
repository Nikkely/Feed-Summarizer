package rss

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
)

// FeedFetcher defines an interface for fetching RSS feeds.
type FeedFetcher interface {
	Fetch(feedUrl string) (*gofeed.Feed, error)
}

// GoFeedFetcher is a concrete implementation of FeedFetcher using gofeed.
type GoFeedFetcher struct{}

// Fetch fetches the RSS feed using gofeed.
func (g *GoFeedFetcher) Fetch(feedUrl string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	return fp.ParseURL(feedUrl)
}

// FetchRss fetches RSS feed items from the given feed URL using the provided FeedFetcher.
// It returns the feed items as a JSON string and an error if any issues occur.
//
// Parameters:
//   - fetcher: An implementation of the FeedFetcher interface.
//   - feedUrl: A string representing the URL of the RSS feed.
//
// Returns:
//   - []string: A slice of strings containing the fetched RSS feed json.
//   - error: An error object if an error occurs, otherwise nil.
func fetchRss(fetcher FeedFetcher, feedUrl string) (string, error) {
	feed, err := fetcher.Fetch(feedUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed from URL %s: %w", feedUrl, err)
	}

	jsonData, err := json.Marshal(feed.Items)
	if err != nil {
		return "", fmt.Errorf("failed to marshal RSS feed items: %w", err)
	}

	return string(jsonData), nil
}

// FetchRss is a convenience function that uses the GoFeedFetcher to fetch RSS feeds.
// It wraps the fetchRss function with a default FeedFetcher implementation.
func FetchRss(feedUrl string) (string, error) {
	return fetchRss(&GoFeedFetcher{}, feedUrl)
}
