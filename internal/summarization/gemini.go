package summarization

import (
	"context"

	"google.golang.org/genai"
)

// GeminiClient is a client for interacting with the Gemini API.
// It allows setting API keys, configuring endpoints, and summarizing page URLs.
type GeminiClient struct {
	model string
}

func NewGeminiClient(model string) *GeminiClient {
	return &GeminiClient{
		model: model,
	}
}

// Summarize summarizes the content of the given page URLs using the Gemini API.
// Parameters:
//   - pageUrl: A slice of strings representing the URLs to summarize.
//
// Returns:
//   - string: The summarized content.
//   - error: An error if the summarization fails.
func (g *GeminiClient) Send(prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
