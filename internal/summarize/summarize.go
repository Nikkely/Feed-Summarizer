package summarize

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"rss-summarizer/pkg/prompt"
	"rss-summarizer/pkg/rss"
)

// RSSInfo represents the title, link, and optional page content of an RSS feed item.
type RSSInfo struct {
	// Title is the title of the RSS feed item.
	Title string `json:"title"`

	// Link is the URL of the RSS feed item.
	Link string `json:"link"`

	Page string `json:"page,omitempty"`
}

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
//   - feedURL: A string representing the URL of the RSS feed.
//
// Returns:
//   - string: The generated summary.
//   - error: An error if any of the URLs cannot be processed or summarization fails.
func Summarize(client GenAIClient, feedURL string) (string, error) {
	var err error
	res, err := rss.FetchRSS(feedURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed from URL %s: %w", feedURL, err)
	}

	if len(res) == 0 {
		return "", fmt.Errorf("no RSS info provided for summarization")
	}

	var infos []RSSInfo
	for _, item := range res {
		page, htmlErr := fetchHTML(item.Link)
		if htmlErr != nil {
			err = fmt.Errorf("%w failed to fetch HTML for URL %s: %w;", err, item.Link, htmlErr)
		} else {
			infos = append(infos, RSSInfo{
				Title: item.Title,
				Link:  item.Link,
				Page:  page,
			})
		}
	}
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

	hoge := promptBuilder.Build()
	return client.Send(hoge)
}
