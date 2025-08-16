package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// httpClientTimeout defines the timeout duration for HTTP client requests.
	httpClientTimeout = 60 * time.Second

	// contextTimeout defines the timeout duration for the context used in FetchHTMLPages.
	contextTimeout = 3 * time.Minute

	// maxConcurrentGoroutines defines the maximum number of concurrent goroutines allowed in FetchHTMLPages.
	maxConcurrentGoroutines = 10
)

// HTMLPageFetcher defines a function type for fetching HTML content of a given URL.
// Parameters:
//   - string: The URL of the page to fetch.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the fetch operation fails.
type HTMLPageFetcher func(string) (string, error)

// FetchHTML retrieves the HTML content of the given URL as a string.
// Parameters:
//   - url: A string representing the target URL.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the request or reading the response fails.
func FetchHTML(url string) (html string, err error) {
	c := &http.Client{
		Timeout: httpClientTimeout,
	}
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeError := resp.Body.Close(); closeError != nil {
			err = errors.Join(err, fmt.Errorf("error closing response body: %w", closeError))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// FetchHTMLPages fetches the HTML content for multiple URLs concurrently.
// It uses the provided HTMLPageFetcher function to retrieve the content of each URL.
// If a timeout occurs or an error happens during fetching, the error is logged and processing continues.
//
// Parameters:
//   - urls: A slice of strings representing the URLs to fetch.
//   - fetcher: A function that fetches the HTML content of a given URL.
//
// Returns:
//   - map[string]string: A map where the keys are URLs and the values are their corresponding HTML content.
//   - error: An aggregated error if any of the URLs cannot be processed.
func FetchHTMLPages(urls []string, fetcher HTMLPageFetcher) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error
	result := make(map[string]string)
	semaphore := make(chan struct{}, maxConcurrentGoroutines)

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(url string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				mu.Lock()
				err = errors.Join(err, fmt.Errorf("context timeout; fetching URL: %s", url))
				mu.Unlock()
			default:
				page, htmlErr := fetcher(url)
				mu.Lock()
				if htmlErr != nil {
					err = errors.Join(err, htmlErr)
				} else {
					result[url] = page
				}
				mu.Unlock()
			}
		}(url)
	}

	wg.Wait()

	if err != nil {
		err = errors.Join(err, fmt.Errorf("some URLs may not have been fetched successfully"))
	}
	return result, err
}
