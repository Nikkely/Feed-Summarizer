package summarize

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockGenAIClient struct{}

func (m *MockGenAIClient) Send(prompt string) (string, error) {
	if prompt == "error" {
		return "", errors.New("mock error")
	}
	return "mock summary", nil
}

func TestSummarize(t *testing.T) {
	mockClient := &MockGenAIClient{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>
			<rss version="2.0">
				<channel>
					<title>Test Feed</title>
					<item>
						<title>Test Item</title>
						<link>http://example.com/test</link>
						<description>Test Description</description>
					</item>
				</channel>
			</rss>`))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	result, err := Summarize(mockClient, ts.URL)
	assert.NoError(t, err, "Summarize returned an unexpected error")
	assert.Equal(t, "mock summary", result, "Summarize result mismatch")
}

func TestSummarize_NoURLs(t *testing.T) {
	mockClient := &MockGenAIClient{}
	_, err := Summarize(mockClient, "")
	assert.Error(t, err, "Expected an error but got nil")
}

func TestSummarize_FetchHTMLError(t *testing.T) {
	mockClient := &MockGenAIClient{}
	_, err := Summarize(mockClient, "http://invalid-url")
	assert.Error(t, err, "Expected an error but got nil")
}

func TestSummarize_EmptyFeed(t *testing.T) {
	mockClient := &MockGenAIClient{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>
			<rss version="2.0">
				<channel>
					<title>Empty Feed</title>
				</channel>
			</rss>`))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, err := Summarize(mockClient, ts.URL)
	assert.Error(t, err, "Expected an error for empty feed but got nil")
}

func TestSummarize_MultipleItems(t *testing.T) {
	mockClient := &MockGenAIClient{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>
			<rss version="2.0">
				<channel>
					<title>Test Feed</title>
					<item>
						<title>Item 1</title>
						<link>http://example.com/item1</link>
						<description>Description 1</description>
					</item>
					<item>
						<title>Item 2</title>
						<link>http://example.com/item2</link>
						<description>Description 2</description>
					</item>
				</channel>
			</rss>`))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	result, err := Summarize(mockClient, ts.URL)
	assert.NoError(t, err, "Summarize returned an unexpected error")
	assert.Equal(t, "mock summary", result, "Summarize result mismatch")
}

func TestSummarize_InvalidRSS(t *testing.T) {
	mockClient := &MockGenAIClient{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid RSS Content"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, err := Summarize(mockClient, ts.URL)
	assert.Error(t, err, "Expected an error for invalid RSS but got nil")
}
