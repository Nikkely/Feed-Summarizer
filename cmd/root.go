// Package cmd implements command-line interface for RSS feed summarizer.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Command line flags
var (
	// genAPIKind specifies the generative AI API to use for summarization
	genAPIKind string
	// systemPromptPath is the path to custom system prompt template file
	systemPromptPath string
	// userPromptPath is the path to custom user prompt template file
	userPromptPath string
	// formatOutput determines whether to format the output as JSON
	formatOutput bool
	// outputTemplatePath is the path to custom output template file
	outputTemplatePath string
	// outputDest specifies where to send the output (standard, file, or datastore)
	outputDest string
)

// rootCmd is the base command for RSS feed summarizer CLI.
// It accepts a URL argument and supports various flags for customization.
var rootCmd = &cobra.Command{
	Use:   "summarize [url]",
	Short: "Summarize RSS feed content using AI",
	Long: `A CLI tool that fetches RSS feed content and generates summaries using AI.
It supports custom prompts and output formatting.

Example:
  summarize https://example.com/feed.xml
  summarize https://example.com/feed.xml --format`,
	Args: cobra.MinimumNArgs(1),
	RunE: summarize,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&genAPIKind, "gen-api-kind", "gemini", "Generative AI API type (currently only 'gemini' is supported)")
	rootCmd.Flags().StringVar(&systemPromptPath, "system-prompt", "", "Path to custom system prompt template file")
	rootCmd.Flags().StringVar(&userPromptPath, "user-prompt", "", "Path to custom user prompt template file")
	rootCmd.Flags().BoolVar(&formatOutput, "format", false, "Format output as JSON with template")
	rootCmd.Flags().StringVar(&outputTemplatePath, "output-template", "", "Custom output template path (only used when -format is true)")
	// rootCmd.Flags().StringVar(&outputDest, "output-dest", "standard", "Output destination (e.g., 'standard', 'file', 'datastore')")
}
