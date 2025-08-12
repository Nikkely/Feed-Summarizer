package prompt

import "fmt"

// PromptBuilder is a utility for constructing prompts for summarization.
// It allows appending user input to a system-defined prompt format.
type PromptBuilder struct {
	// SystemPrompt is the base prompt provided by the system.
	SystemPrompt string

	// UserPromptFmt is the format string for user prompts, compatible with fmt.Printf.
	UserPromptFmt string

	// userPrompt stores the accumulated user input.
	userPrompt string
}

// NewPromptBuilder creates a new instance of PromptBuilder.
// Parameters:
//   - systemPrompt: A string representing the base system prompt.
//   - userPromptFmt: A string representing the format for user prompts.
//
// Returns:
//   - *PromptBuilder: A pointer to the newly created PromptBuilder instance.
func NewPromptBuilder(systemPrompt string, userPromptFmt string) *PromptBuilder {
	return &PromptBuilder{
		SystemPrompt:  systemPrompt,
		UserPromptFmt: userPromptFmt,
	}
}

// Append appends user input to the prompt using the specified format.
// Parameters:
//   - text: A variadic slice of strings representing the user input to append.
//
// Returns:
//   - *PromptBuilder: The updated PromptBuilder instance.
func (p *PromptBuilder) Append(text ...string) *PromptBuilder {
	args := make([]any, len(text))
	for i, v := range text {
		args[i] = v
	}
	p.userPrompt += fmt.Sprintf(p.UserPromptFmt, args...)
	return p
}

// Build constructs the final prompt by combining the system prompt and user input.
// Returns:
//   - string: The complete prompt ready for use.
func (p *PromptBuilder) Build() string {
	return p.SystemPrompt + p.userPrompt
}
