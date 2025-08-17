// Package summarize provides functionality for summarizing RSS feed content.
//
// This package offers the following key features:
//   - Fetching and parsing RSS feeds
//   - Retrieving HTML content for feed articles
//   - Generating summaries using AI-powered content analysis
//   - Customizable template-based output formatting
//
// Use the Summarizer to generate summaries from RSS feed URLs.
package summarize

import (
	"fmt"
	"log"
	"os"
	"text/template"

	genAi "feed-summarizer/internal/ai_client"
	"feed-summarizer/internal/fetcher"
	"feed-summarizer/pkg/prompt"

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
// It fetches HTML content for each feed item using the provided pageFetcher.
// If an error occurs while fetching a page, it logs the error and continues processing other items.
// Parameters:
//   - feed: A pointer to a gofeed.Feed object containing RSS feed data.
//   - pageFetcher: A function that fetches the HTML content of a given URL.
//
// Returns:
//   - []RSSInfo: A slice of RSSInfo containing the title, link, and optional page content.
//   - error: An aggregated error if any of the URLs cannot be processed.
func NewRSSInfo(feed *gofeed.Feed, pageFetcher fetcher.HTMLPageFetcher) (infos []RSSInfo, err error) {
	urls := make([]string, 0, len(feed.Items))
	for _, item := range feed.Items {
		urls = append(urls, item.Link)
	}
	pages, err := fetcher.FetchHTMLPages(urls, pageFetcher)

	for _, item := range feed.Items {
		page := pages[item.Link]
		infos = append(infos, RSSInfo{
			Title: item.Title,
			Link:  item.Link,
			Page:  page, // This value will be nil if the retrieval fails.
		})
	}
	return infos, err
}

// Summarizer is responsible for summarizing RSS feed content using a GenAIClient.
// It fetches RSS feeds, retrieves HTML content, and generates summaries based on prompts.
//
// The Summarizer workflow:
// 1. Fetches RSS feed from the provided URL
// 2. Extracts feed items and their HTML content
// 3. Builds prompts using customizable templates
// 4. Generates summaries using AI
type Summarizer struct {
	client        genAi.GenAIClient
	feedFetcher   fetcher.FeedFetcher
	pageFetcher   fetcher.HTMLPageFetcher
	promptBuilder *prompt.PromptBuilder
}

// NewSummarizer initializes a new Summarizer instance.
// Parameters:
//   - client: An instance of GenAIClient for generating summaries.
//   - feedFetcher: A function to fetch RSS feeds.
//   - pageFetcher: A function to fetch HTML content of URLs.
//
// Returns:
//   - *Summarizer: A new Summarizer instance.
func NewSummarizer(client genAi.GenAIClient, feedFetcher fetcher.FeedFetcher, pageFetcher fetcher.HTMLPageFetcher) *Summarizer {
	promptBuilder := prompt.NewPromptBuilder(systemPrompt, userPromptTemplate)
	return &Summarizer{
		client:        client,
		feedFetcher:   feedFetcher,
		pageFetcher:   pageFetcher,
		promptBuilder: promptBuilder,
	}
}

// LoadPromptBuilder initializes the prompt builder with system and user prompts.
// Parameters:
//   - sysPromptTxtPath: Path to the system prompt text file.
//   - usrPromptTmplPath: Path to the user prompt template file.
//
// Returns:
//   - error: An error if loading or parsing the files fails.
func (s *Summarizer) LoadPromptBuilder(sysPromptTxtPath, usrPromptTmplPath string) error {
	sysPrompt, err := txtFileLoader(sysPromptTxtPath)
	if err != nil {
		return err
	}

	usrPromptTmpl, err := template.ParseFiles(usrPromptTmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", usrPromptTmplPath, err)
	}

	s.promptBuilder = prompt.NewPromptBuilder(sysPrompt, usrPromptTmpl)
	return nil
}

// Summarize generates a summary for the content of the given RSS feed URL.
// It continues processing even if some HTML pages fail to fetch, logging the errors.
// Parameters:
//   - feedURL: A string representing the URL of the RSS feed.
//
// Returns:
//   - string: The generated summary.
//   - error: An error if the summarization process fails entirely.
func (s *Summarizer) Summarize(feedURL string) (string, error) {
	var err error
	if s.promptBuilder == nil {
		return "", fmt.Errorf("prompt builder is not initialized")
	}

	feed, err := s.feedFetcher(feedURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	infos, err := NewRSSInfo(feed, s.pageFetcher)
	if err != nil {
		log.Printf("failed to fetch HTML for some URLs: %v", err) // Continue if page retrieval fails
	}

	for _, info := range infos {
		s.promptBuilder.Append(info)
	}

	return s.client.Send(s.promptBuilder.Build())
}

// txtFileLoader reads the content of a text file and returns it as a string.
// Parameters:
//   - filePath: A string representing the path to the text file.
//
// Returns:
//   - string: The content of the file as a string.
//   - error: An error if the file cannot be read.
func txtFileLoader(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(content), nil
}
