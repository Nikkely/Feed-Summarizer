// Package fetcher provides functionality for fetching and processing external resources.
// It includes methods for fetching RSS feeds and HTML pages, such as FeedFetcher and HTMLPageFetcher.
// The package is designed with dependency injection in mind, enabling easier testing of components
// that rely on external data fetching.
package fetcher

import (
	"errors"
	"fmt"

	"github.com/mmcdole/gofeed"
)

// FeedFetcher defines a function type for fetching an RSS feed.
// Parameters:
//   - string: The URL of the RSS feed to fetch.
//
// Returns:
//   - *gofeed.Feed: The parsed RSS feed.
//   - error: An error if the fetch operation fails.
type FeedFetcher func(string) (*gofeed.Feed, error)

// FetchFeed fetches and parses an RSS feed from the given URL.
// Parameters:
//   - feedURL: A string representing the URL of the RSS feed.
//
// Returns:
//   - *gofeed.Feed: The parsed RSS feed.
//   - error: An error if the fetch operation fails.
func FetchFeed(feedURL string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("failed to fetch RSS feed from URL %s: %w", feedURL, err))
	}
	return feed, err
}
