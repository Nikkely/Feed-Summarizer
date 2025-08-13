package summarize

import (
	"fmt"

	"rss-summarizer/pkg/prompt"
	"rss-summarizer/pkg/rss"
)

// Summarize generates a summary for the content of the given URLs using a GenAIClient.
// Parameters:
//   - client: An instance of GenAIClient to handle the summarization.
//   - feedURL: A string representing the URL of the RSS feed.
//
// Returns:
//   - string: The generated summary.
//   - error: An error if any of the URLs cannot be processed or summarization fails.
func Summarize(client GenAIClient, feedURL string) (string, error) {
	infos, err := rss.FetchRSS(feedURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed from URL %s: %w", feedURL, err)
	}

	if len(infos) == 0 {
		return "", fmt.Errorf("no RSS info provided for summarization")
	}

	promptBuilder, err := prompt.NewPromptBuilder(`次のページを要約してください`, `
概要：{{.title}}, URL:{{.link}} 
{{ with .Page }}
  {{ . }}
{{ end }}
	`)
	if err != nil {
		return "", fmt.Errorf("failed to create prompt builder: %w", err)
	}

	for _, info := range infos {
		promptBuilder.Append(info)
	}

	return client.Send(promptBuilder.Build())
}
