package cmd

import (
	genAi "feed-summarizer/internal/ai_client"
	"feed-summarizer/internal/fetcher"
	"feed-summarizer/internal/jsonify"
	sum "feed-summarizer/internal/summarize"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
)

func summarize(cmd *cobra.Command, args []string) error {
	sumClient := genAi.NewGenAIClient(genAPIKind)
	if sumClient == nil {
		return fmt.Errorf("unsupported API type: %s", genAPIKind)
	}

	summarizer := sum.NewSummarizer(sumClient, fetcher.FetchFeed, fetcher.FetchHTML)
	if systemPromptPath != "" && userPromptPath != "" {
		if err := summarizer.LoadPromptBuilder(systemPromptPath, userPromptPath); err != nil {
			return fmt.Errorf("failed to load prompt builder: %w", err)
		}
	}

	for _, url := range args {
		summary, err := summarizer.Summarize(url)
		if err != nil {
			return fmt.Errorf("failed to summarize feed: %w", err)
		}

		if !formatOutput {
			fmt.Println(summary)
			return nil
		}

		var outputTemplate *template.Template
		if outputTemplatePath != "" {
			outputTemplate, err = template.ParseFiles(outputTemplatePath)
			if err != nil {
				return fmt.Errorf("failed to read output template from %s: %w", outputTemplatePath, err)
			}
		} else {
			outputTemplate = jsonify.OutputTemplate
		}

		formattedResult, err := jsonify.ExtractAndFormat(summary, outputTemplate)
		if err != nil {
			return fmt.Errorf("failed to format summary: %w", err)
		}

		switch outputDest {
		default:
			fmt.Println(formattedResult)
		}
	}

	return nil
}
