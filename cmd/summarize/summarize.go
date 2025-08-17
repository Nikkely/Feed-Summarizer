package main

import (
	"flag"
	"fmt"
	"log"
	genAi "feed-summarizer/internal/ai_client"
	"feed-summarizer/internal/fetcher"
	"feed-summarizer/internal/jsonify"
	sum "feed-summarizer/internal/summarize"
	"text/template"
)

func main() {
	// Set up command-line flags
	urlArg := flag.String("url", "https://example.com/feed", "RSS feed URL to summarize")
	genAPIKindArg := flag.String("gen-api-kind", "gemini", "Generative AI API type (currently only 'gemini' is supported)")
	systemPromptPathArg := flag.String("system-prompt", "", "Path to custom system prompt template file")
	userPromptPathArg := flag.String("user-prompt", "", "Path to custom user prompt template file")
	formatArg := flag.Bool("format", false, "Format output as JSON with template")
	outputTemplatePathArg := flag.String("output", "", "Custom output template path (only used when -format is true)")
	flag.Parse()

	// Validate command-line arguments
	if *urlArg == "" || *urlArg == "https://example.com/feed" {
		log.Fatal("Please provide a valid RSS feed URL")
	}

	sumClient := genAi.NewGenAIClient(*genAPIKindArg)
	if sumClient == nil {
		log.Fatal("Unsupported API")
	}

	summarizer := sum.NewSummarizer(sumClient, fetcher.FetchFeed, fetcher.FetchHTML)
	if systemPromptPathArg != nil && *systemPromptPathArg != "" && userPromptPathArg != nil && *userPromptPathArg != "" {
		if err := summarizer.LoadPromptBuilder(*systemPromptPathArg, *userPromptPathArg); err != nil {
			log.Fatal("Failed to load prompt builder:", err)
		}
	}
	summary, err := summarizer.Summarize(*urlArg)
	if err != nil {
		log.Fatal("Failed to summarize Feed:", err)
	}

	if !*formatArg {
		// Output raw summary when formatting is not requested
		fmt.Println(summary)
		return
	}

	// Set up output template for formatting
	var outputTemplate *template.Template
	if *outputTemplatePathArg != "" {
		// Load custom template if specified
		outputTemplate, err = template.ParseFiles(*outputTemplatePathArg)
		if err != nil {
			log.Fatalf("Failed to read output template from %s: %v", *outputTemplatePathArg, err)
		}
	} else {
		// Use default template from jsonify package
		outputTemplate = jsonify.OutputTemplate
	}

	// Format the summary using the template
	formattedResult, err := jsonify.ExtractAndFormat(summary, outputTemplate)
	if err != nil {
		log.Fatalf("Failed to format summary: %v", err)
	}

	fmt.Println(formattedResult)
}
