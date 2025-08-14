package summarize

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"rss-summarizer/pkg/prompt"

	"github.com/mmcdole/gofeed"
)

// HTMLPageFetcher defines a function type for fetching HTML content of a given URL.
// Parameters:
//   - string: The URL of the page to fetch.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the fetch operation fails.
type HTMLPageFetcher func(string) (string, error)

// FeedFetcher defines a function type for fetching an RSS feed.
// Parameters:
//   - string: The URL of the RSS feed to fetch.
//
// Returns:
//   - *gofeed.Feed: The parsed RSS feed.
//   - error: An error if the fetch operation fails.
type FeedFetcher func(string) (*gofeed.Feed, error)

// RSSInfo represents the title, link, and optional page content of an RSS feed item.
type RSSInfo struct {
	// Title is the title of the RSS feed item.
	Title string `json:"title"`

	// Link is the URL of the RSS feed item.
	Link string `json:"link"`

	// Page contains the optional HTML content of the RSS feed item.
	// It is omitted from the JSON output if empty.
	Page string `json:"page,omitempty"`
}

// NewRSSInfo creates a slice of RSSInfo from a gofeed.Feed.
// Parameters:
//   - feed: A pointer to a gofeed.Feed object containing RSS feed data.
//   - pageFetcher: A function that fetches the HTML content of a given URL.
//
// Returns:
//   - []RSSInfo: A slice of RSSInfo containing the title, link, and optional page content.
//   - error: An error if any of the URLs cannot be processed.
func NewRSSInfo(feed *gofeed.Feed, pageFetcher HTMLPageFetcher) (infos []RSSInfo, err error) {
	for _, item := range feed.Items {
		page, htmlErr := pageFetcher(item.Link) // TODO: マルチスレッド実行可能に
		if htmlErr != nil {
			err = errors.Join(err, htmlErr)
			continue
		}

		infos = append(infos, RSSInfo{
			Title: item.Title,
			Link:  item.Link,
			Page:  page,
		})
	}
	if err != nil {
		err = errors.Join(err, fmt.Errorf("failed to fetch HTML for some URLs: %w", err))
	}
	return infos, err
}

// FetchHTML retrieves the HTML content of the given URL as a string.
// Parameters:
//   - url: A string representing the target URL.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the request or reading the response fails.
func FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

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

// Summarize generates a summary for the content of the given URLs using a GenAIClient.
// Parameters:
//   - client: An instance of GenAIClient to handle the summarization.
//   - feedFetcher: A function that fetches an RSS feed from a given URL.
//   - pageFetcher: A function that fetches the HTML content of a given URL.
//   - feedURL: A string representing the URL of the RSS feed.
//
// Returns:
//   - string: The generated summary.
//   - error: An error if any of the URLs cannot be processed or summarization fails.
//
// Example:
//   client := NewGenAIClient()
//   summary, err := Summarize(client, FetchFeed, FetchHTML, "https://example.com/rss")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(summary)
func Summarize(client GenAIClient, feedFetcher FeedFetcher, pageFetcher HTMLPageFetcher, feedURL string) (string, error) {
	var err error
	feed, err := feedFetcher(feedURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	infos, err := NewRSSInfo(feed, pageFetcher)
	if err != nil {
		return "", fmt.Errorf("failed to fetch HTML for some URLs: %w", err)
	}

	temp, err := template.New("info").Parse(`
概要：{{.Title}}, URL:{{.Link}} 
{{ with .Page }}
  {{ . }}
{{ end }}
	`)
	if err != nil {
		return "", fmt.Errorf("failed to parse user prompt template: %w", err)
	}

	promptBuilder := prompt.NewPromptBuilder(`次のページを要約してください`, temp)
	for _, info := range infos {
		promptBuilder.Append(info)
	}

	return client.Send(promptBuilder.Build())
}
