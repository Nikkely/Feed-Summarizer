package summarize

import (
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestFetchHTML(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test Page</body></html>"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	result, err := fetchHTML(ts.URL)
	assert.NoError(t, err, "fetchHTML returned an unexpected error")
	assert.Contains(t, result, "Test Page", "fetchHTML result mismatch")
}

func TestFetchHTML_Error(t *testing.T) {
	_, err := fetchHTML("http://invalid-url")
	assert.Error(t, err, "Expected an error but got nil")
}

func TestNewRSSInfo(t *testing.T) {
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
		return "", errors.New("failed to fetch page")
	}

	infos, err := NewRSSInfo(mockFeed, mockPageFetcher)
	assert.Error(t, err, "NewRSSInfo returned an unexpected error")
	assert.Len(t, infos, 1, "Expected 1 valid RSSInfo but got a different count")
	assert.Equal(t, "Item 1", infos[0].Title, "RSSInfo title mismatch")
	assert.Equal(t, "http://example.com/item1", infos[0].Link, "RSSInfo link mismatch")
	assert.Contains(t, infos[0].Page, "Page 1", "RSSInfo page content mismatch")
}
