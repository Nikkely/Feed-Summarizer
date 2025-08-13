package rss

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
)

// RSSInfo represents the title and link of an RSS feed item.
type RSSInfo struct {
	// Title is the title of the RSS feed item.
	Title string `json:"title"`

	// Link is the URL of the RSS feed item.
	Link string `json:"link"`
}

// FeedFetcher defines an interface for fetching RSS feeds.
// It provides a method to fetch RSS feed items from a given URL.
type FeedFetcher interface {
	// Fetch retrieves RSS feed items from the specified URL.
	// Parameters:
	//   - feedURL: A string representing the URL of the RSS feed.
	// Returns:
	//   - []RSSInfo: A slice of RSSInfo containing the fetched feed items.
	//   - error: An error if the fetch operation fails.
	Fetch(feedURL string) ([]RSSInfo, error)
}

// GoFeedFetcher is a concrete implementation of FeedFetcher using gofeed.
type GoFeedFetcher struct{}

// Fetch fetches the RSS feed using gofeed.
// Parameters:
//   - feedURL: A string representing the URL of the RSS feed.
// Returns:
//   - []RSSInfo: A slice of RSSInfo containing the fetched feed items.
//   - error: An error if the fetch operation fails.
func (g *GoFeedFetcher) Fetch(feedURL string) ([]RSSInfo, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}
	return g.translate(feed)
}

// translate converts a gofeed.Feed into a slice of RSSInfo.
// Parameters:
//   - feed: A pointer to a gofeed.Feed object.
// Returns:
//   - []RSSInfo: A slice of RSSInfo containing the feed items.
//   - error: An error if the feed is nil.
func (g *GoFeedFetcher) translate(feed *gofeed.Feed) ([]RSSInfo, error) {
	if feed == nil {
		return nil, fmt.Errorf("feed is nil")
	}

	var info []RSSInfo
	for _, item := range feed.Items {
		info = append(info, RSSInfo{
			Title: item.Title,
			Link:  item.Link,
		})
	}
	return info, nil
}

// fetchRSS fetches RSS feed items from the given feed URL using the provided FeedFetcher.
// Parameters:
//   - fetcher: An implementation of the FeedFetcher interface.
//   - feedURL: A string representing the URL of the RSS feed.
// Returns:
//   - []RSSInfo: A slice of RSSInfo containing the fetched feed items.
//   - error: An error if the fetch operation fails.
func fetchRSS(fetcher FeedFetcher, feedURL string) ([]RSSInfo, error) {
	info, err := fetcher.Fetch(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed from URL %s: %w", feedURL, err)
	}
	return info, nil
}

// FetchRSStoJSONString fetches RSS feed items and returns them as a JSON string.
// Parameters:
//   - feedURL: A string representing the URL of the RSS feed.
// Returns:
//   - string: A JSON string containing the fetched feed items.
//   - error: An error if the fetch or JSON marshaling operation fails.
func FetchRSStoJSONString(feedURL string) (string, error) {
	info, err := fetchRSS(&GoFeedFetcher{}, feedURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	jsonData, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("failed to marshal RSS feed items: %w", err)
	}

	return string(jsonData), nil
}

// FetchRSS fetches RSS feed items and returns them as a slice of RSSInfo.
// Parameters:
//   - feedURL: A string representing the URL of the RSS feed.
// Returns:
//   - []RSSInfo: A slice of RSSInfo containing the fetched feed items.
//   - error: An error if the fetch operation fails.
func FetchRSS(feedURL string) ([]RSSInfo, error) {
	info, err := fetchRSS(&GoFeedFetcher{}, feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	return info, nil
}
