package fetcher

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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
	resp, err := http.Get(url)
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
