package summarize

import (
	"errors"
	"fmt"
	"text/template"

	"rss-summarizer/pkg/prompt"

	"github.com/mmcdole/gofeed"
)

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

// Summarize generates a summary for the content of the given URLs using a GenAIClient.
// Parameters:
//   - client: An instance of GenAIClient to handle the summarization.
//   - feedURL: A string representing the URL of the RSS feed.
//
// Returns:
//   - string: The generated summary.
//   - error: An error if any of the URLs cannot be processed or summarization fails.
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
