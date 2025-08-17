// Package aiclient provides an interface for AI-based summarization clients.
// It defines the GenAIClient interface, which is used to send text to an AI model
// and receive a summarized response.
package aiclient

// Client defines an interface for summarization clients.
type GenAIClient interface {
	// Summarize summarizes the content of the given page URLs.
	// Parameters:
	//   - pageUrl: A slice of strings representing the URLs to summarize.
	// Returns:
	//   - string: The summarized content.
	//   - error: An error if the summarization fails.
	Send(text string) (string, error)
}

// NewGenAIClient creates a new AI client of the specified type.
// Parameters:
//   - kind: Type of AI client to create. Currently only "gemini" is supported.
//
// Returns:
//   - GenAIClient: A new instance of the AI client, or nil if the type is not supported.
func NewGenAIClient(kind string) GenAIClient {
	switch kind {
	case "gemini":
		return NewGeminiClient("gemini-2.5-flash-lite")
	default:
		return nil
	}
}
