package fetcher

import (
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
		// TODO: 処理を止めずにエラーログを出力する
		return nil, fmt.Errorf("failed to fetch RSS feed from URL %s: %w", feedURL, err)
	}
	return feed, nil
}
