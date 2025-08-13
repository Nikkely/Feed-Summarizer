package summarize

import (
	"fmt"
	"io"
	"net/http"

	"rss-summarizer/pkg/prompt"
)

// fetchHTML retrieves the HTML content of the given URL as a string.
// Parameters:
//   - url: A string representing the target URL.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the request or reading the response fails.
func fetchHTML(url string) (string, error) {
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

// Summarize generates a summary for the content of the given URLs using a GenAIClient.
// Parameters:
//   - client: An instance of GenAIClient to handle the summarization.
//   - urls: A slice of strings representing the URLs to summarize.
//
// Returns:
//   - string: The generated summary.
//   - error: An error if any of the URLs cannot be processed or summarization fails.
func Summarize(client GenAIClient, urls []string) (string, error) {
	if len(urls) == 0 {
		return "", fmt.Errorf("no URLs provided for summarization")
	}

	promptBuilder := prompt.NewPromptBuilder("次のページを要約してください", "ページ:%s")
	for _, url := range urls {
		content, err := fetchHTML(url)
		if err != nil {
			return "", fmt.Errorf("error fetching HTML from %s: %w", url, err)
		}
		promptBuilder.Append(content)
	}

	return client.Send(promptBuilder.Build())
}
