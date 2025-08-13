package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

// PromptBuilder is a utility for constructing prompts for summarization.
// It allows appending user input to a system-defined prompt format using templates.
type PromptBuilder struct {
	// SystemPrompt is the base prompt provided by the system.
	SystemPrompt string

	// UserPromptTemplate is the template for user prompts.
	UserPromptTemplate *template.Template

	// userPrompt stores the accumulated user input.
	userPrompt string
}

// NewPromptBuilder creates a new instance of PromptBuilder.
// Parameters:
//   - systemPrompt: A string representing the base system prompt.
//   - userPromptTemplate: A string representing the template for user prompts.
//
// Returns:
//   - *PromptBuilder: A pointer to the newly created PromptBuilder instance.
//   - error: An error if the template parsing fails.
func NewPromptBuilder(systemPrompt string, userPromptTemplate string) (*PromptBuilder, error) {
	temp, err := template.New("info").Parse(userPromptTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user prompt template: %w", err)
	}

	return &PromptBuilder{
		SystemPrompt:       systemPrompt,
		UserPromptTemplate: temp,
	}, nil
}

// Append appends user input to the prompt using the specified template.
// Parameters:
//   - vars: A map or struct containing keys and values to populate the template.
//
// Returns:
//   - *PromptBuilder: The updated PromptBuilder instance, or nil if an error occurs.
func (p *PromptBuilder) Append(vars any) *PromptBuilder {
	var buf bytes.Buffer
	if err := p.UserPromptTemplate.Execute(&buf, vars); err != nil {
		return nil
	}
	p.userPrompt += buf.String()

	return p
}

// Build constructs the final prompt by combining the system prompt and user input.
// Returns:
//   - string: The complete prompt ready for use.
func (p *PromptBuilder) Build() string {
	return p.SystemPrompt + p.userPrompt
}
