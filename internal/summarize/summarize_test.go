package summarize

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"rss-summarizer/internal/fetcher"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

type MockGenAIClient struct{}

func (m *MockGenAIClient) Send(prompt string) (string, error) {
	if prompt == "error" {
		return "", errors.New("mock error")
	}
	return "mock summary", nil
}

func TestSummarize_Updated(t *testing.T) {
	mockClient := &MockGenAIClient{}
	mockFeedFetcher := func(_ string) (*gofeed.Feed, error) {
		return &gofeed.Feed{
			Items: []*gofeed.Item{
				{Title: "Test Item", Link: "http://example.com/test"},
			},
		}, nil
	}
	mockPageFetcher := func(_ string) (string, error) {
		return "<html>Test Page</html>", nil
	}

	s := NewSummarizer(mockClient, mockFeedFetcher, mockPageFetcher)
	if err := s.LoadPromptBuilder("../../templates/system_prompt.txt", "../../templates/user_prompt.tmpl"); err != nil {
		log.Fatalln(err)
	}
	result, err := s.Summarize("http://example.com/rss")
	assert.NoError(t, err, "Summarize returned an unexpected error")
	assert.Equal(t, "mock summary", result, "Summarize result mismatch")
}

func TestFetchFeed(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>
			<rss version="2.0">
				<channel>
					<title>Test Feed</title>
					<item>
						<title>Test Item</title>
						<link>http://example.com/test</link>
					</item>
				</channel>
			</rss>`)); err != nil {
			log.Fatalf("Error writing response: %v", err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	feed, err := fetcher.FetchFeed(ts.URL)
	assert.NoError(t, err, "FetchFeed returned an unexpected error")
	assert.NotNil(t, feed, "Expected a valid feed but got nil")
	assert.Len(t, feed.Items, 1, "Expected 1 item in the feed but got a different count")
	assert.Equal(t, "Test Item", feed.Items[0].Title, "Feed item title mismatch")
}

func TestFetchFeed_Error(t *testing.T) {
	_, err := fetcher.FetchFeed("http://invalid-url")
	assert.Error(t, err, "Expected an error but got nil")
}

func TestLoadPromptBuilder(t *testing.T) {
	s := &Summarizer{}
	err := s.LoadPromptBuilder("../../templates/system_prompt.txt", "../../templates/user_prompt.tmpl")
	assert.NoError(t, err, "LoadPromptBuilder returned an unexpected error")
	assert.NotNil(t, s.promptBuilder, "PromptBuilder should be initialized")
}

func TestSummarize_ErrorWhenPromptBuilderNotInitialized(t *testing.T) {
	mockClient := &MockGenAIClient{}
	mockFeedFetcher := func(_ string) (*gofeed.Feed, error) {
		return &gofeed.Feed{
			Items: []*gofeed.Item{
				{Title: "Test Item", Link: "http://example.com/test"},
			},
		}, nil
	}
	mockPageFetcher := func(_ string) (string, error) {
		return "<html>Test Page</html>", nil
	}

	s := NewSummarizer(mockClient, mockFeedFetcher, mockPageFetcher)
	_, err := s.Summarize("http://example.com/rss")
	assert.Error(t, err, "Expected an error when PromptBuilder is not initialized")
	assert.Contains(t, err.Error(), "prompt builder is not initialized", "Error message mismatch")
}

func TestNewRSSInfo_ErrorHandling(t *testing.T) {
	mockFeed := &gofeed.Feed{
		Items: []*gofeed.Item{
			{Title: "Item 1", Link: "http://example.com/item1"},
			{Title: "Item 2", Link: "http://example.com/item2"},
		},
	}

	mockPageFetcher := func(url string) (string, error) {
		if url == "http://example.com/item1" {
			return "<html>Page 1</html>", nil
		}
		return "", fmt.Errorf("failed to fetch page for URL: %s", url)
	}

	infos, err := NewRSSInfo(mockFeed, mockPageFetcher)
	assert.Error(t, err, "Expected an error but got nil")
	assert.Len(t, infos, 2, "Expected 2 RSSInfo items")
	assert.Equal(t, "Item 1", infos[0].Title, "RSSInfo title mismatch")
	assert.Equal(t, "http://example.com/item1", infos[0].Link, "RSSInfo link mismatch")
	assert.Contains(t, infos[0].Page, "Page 1", "RSSInfo page content mismatch")
	assert.Empty(t, infos[1].Page, "Expected empty page content for failed fetch")
}

func TestNewRSSInfo_WithPageFetcher(t *testing.T) {
	mockFeed := &gofeed.Feed{
		Items: []*gofeed.Item{
			{Title: "Item 1", Link: "http://example.com/item1"},
			{Title: "Item 2", Link: "http://example.com/item2"},
		},
	}

	mockPageFetcher := func(url string) (string, error) {
		if url == "http://example.com/item1" {
			return "<html>Page 1</html>", nil
		}
		return "", fmt.Errorf("failed to fetch page for URL: %s", url)
	}

	infos, err := NewRSSInfo(mockFeed, mockPageFetcher)
	assert.Error(t, err, "Expected an error but got nil")
	assert.Len(t, infos, 2, "Expected 2 RSSInfo items")
	assert.Equal(t, "Item 1", infos[0].Title, "RSSInfo title mismatch")
	assert.Equal(t, "http://example.com/item1", infos[0].Link, "RSSInfo link mismatch")
	assert.Contains(t, infos[0].Page, "Page 1", "RSSInfo page content mismatch")
	assert.Empty(t, infos[1].Page, "Expected empty page content for failed fetch")
}

func TestSummarize_ErrorHandling(t *testing.T) {
	mockClient := &MockGenAIClient{}
	mockFeedFetcher := func(_ string) (*gofeed.Feed, error) {
		return &gofeed.Feed{
			Items: []*gofeed.Item{
				{Title: "Test Item", Link: "http://example.com/test"},
			},
		}, nil
	}
	mockPageFetcher := func(_ string) (string, error) {
		return "", fmt.Errorf("failed to fetch page")
	}

	s := NewSummarizer(mockClient, mockFeedFetcher, mockPageFetcher)
	s.LoadPromptBuilder("../../templates/system_prompt.txt", "../../templates/user_prompt.tmpl")
	result, err := s.Summarize("http://example.com/rss")
	assert.NoError(t, err, "Summarize returned an unexpected error")
	assert.Equal(t, "mock summary", result, "Summarize result mismatch")
}
