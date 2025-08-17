package summarize

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"rss-summarizer/internal/fetcher"
	"rss-summarizer/pkg/prompt"
	"testing"
	"text/template"

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

var testSystemPrompt = `あなたはニュース記事やブログ記事を短く正確にまとめる要約アシスタントです。

# 目的
入力された記事タイトル、URL、およびHTML本文を元に、記事の要点を正確かつ簡潔に日本語でまとめてください。
記事は複数入力されます。なお、記事の主要な事実・数値・固有名詞を落とさず、主観や推測を加えないでください。`

var testUserPromptTemplate = `タイトル：{{.Title}}, URL:{{.Link}} 
{{ with .Page }}
  {{ . }}
{{ end }}`

var testOutputTemplate = `{
  "heading": "{{ .heading }}",
  "summary": "{{ .summary }}"
}`

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

	// テンプレートの直接設定
	s.promptBuilder = prompt.NewPromptBuilder(testSystemPrompt, template.Must(template.New("user").Parse(testUserPromptTemplate)))
	s.outputTmpl = template.Must(template.New("output").Parse(testOutputTemplate))

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

	// 一時ディレクトリにテストデータを書き込んでテスト
	tmpDir := t.TempDir()

	sysPromptPath := filepath.Join(tmpDir, "system_prompt.txt")
	if err := os.WriteFile(sysPromptPath, []byte(testSystemPrompt), 0644); err != nil {
		t.Fatal(err)
	}

	userPromptPath := filepath.Join(tmpDir, "user_prompt.tmpl")
	if err := os.WriteFile(userPromptPath, []byte(testUserPromptTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	err := s.LoadPromptBuilder(sysPromptPath, userPromptPath)
	assert.NoError(t, err, "LoadPromptBuilder returned an unexpected error")
	assert.NotNil(t, s.promptBuilder, "PromptBuilder should be initialized")
	assert.NotNil(t, s.promptBuilder, "PromptBuilder should be initialized")
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
	_ = s.LoadPromptBuilder("../../templates/system_prompt.txt", "../../templates/user_prompt.tmpl")
	result, err := s.Summarize("http://example.com/rss")
	assert.NoError(t, err, "Summarize returned an unexpected error")
	assert.Equal(t, "mock summary", result, "Summarize result mismatch")
}
