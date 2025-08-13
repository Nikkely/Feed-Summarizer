package summarize

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
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
		w.Write([]byte("<html><body>Test Content</body></html>"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	content, err := fetchHTML(ts.URL)
	if err != nil {
		t.Fatalf("fetchHTML returned an error: %v", err)
	}

	expected := "<html><body>Test Content</body></html>"
	if content != expected {
		t.Errorf("expected %q, got %q", expected, content)
	}
}

func TestFetchHTML_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, err := fetchHTML(ts.URL)
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestSummarize(t *testing.T) {
	mockClient := &MockGenAIClient{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test Content"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	urls := []string{ts.URL}
	result, err := Summarize(mockClient, urls)
	if err != nil {
		t.Fatalf("Summarize returned an error: %v", err)
	}

	expected := "mock summary"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSummarize_NoURLs(t *testing.T) {
	mockClient := &MockGenAIClient{}
	_, err := Summarize(mockClient, []string{})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestSummarize_FetchHTMLError(t *testing.T) {
	mockClient := &MockGenAIClient{}
	_, err := Summarize(mockClient, []string{"http://invalid-url"})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}
