package cmd

import (
	"context"
	genAi "feed-summarizer/ai_client"
	db "feed-summarizer/database"
	"feed-summarizer/fetcher"
	"feed-summarizer/jsonify"
	sum "feed-summarizer/summarize"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
)

func summarize(_ *cobra.Command, args []string) error {
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

		formattedResults, err := jsonify.ExtractAndFormat(summary, outputTemplate)
		if err != nil {
			return fmt.Errorf("failed to format summary: %w", err)
		}

		switch outputDest {
		case "datastore":
			ctx := context.Background()
			client, err := db.NewDatastoreClient(ctx, gcpProjectID)
			if err != nil {
				return fmt.Errorf("failed to create datastore client: %w", err)
			}
			keys, err := db.GenerateUUIDs(len(formattedResults))
			if err != nil {
				return fmt.Errorf("failed to generate keys: %w", err)
			}
			if err := client.PutMulti(ctx, "summaries", keys, formattedResults); err != nil {
				return fmt.Errorf("failed to save summary to datastore: %w", err)
			}
		default:
			fmt.Println(formattedResults)
		}
	}

	return nil
}
