package ai_client

import (
	"context"

	"google.golang.org/genai"
)

// GeminiClient is a client for interacting with the Gemini API.
// It allows sending prompts to the Gemini model and retrieving generated content.
type GeminiClient struct {
	// model specifies the Gemini model to use for content generation.
	model string
}

// NewGeminiClient creates a new instance of GeminiClient.
// Parameters:
//   - model: A string representing the Gemini model to use.
//
// Returns:
//   - *GeminiClient: A pointer to the newly created GeminiClient instance.
func NewGeminiClient(model string) *GeminiClient {
	return &GeminiClient{
		model: model,
	}
}

// Send sends a prompt to the Gemini API and retrieves the generated content.
// Parameters:
//   - prompt: A string representing the input prompt to send to the Gemini model.
//
// Returns:
//   - string: The generated content from the Gemini model.
//   - error: An error if the content generation fails.
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
