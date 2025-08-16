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
