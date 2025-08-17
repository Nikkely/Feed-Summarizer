package main

import (
	"flag"
	"fmt"
	"log"
	genAi "rss-summarizer/internal/ai_client"
	"rss-summarizer/internal/fetcher"
	sum "rss-summarizer/internal/summarize"
)

func main() {
	urlArg := flag.String("url", "https://example.com/feed", "rss feed url")
	genApiKindArg := flag.String("gen-api-kind", "gemini", "kind of generative API to use")
	systemPromptPathArg := flag.String("system-prompt", "", "path to system prompt template")
	userPromptPathArg := flag.String("user-prompt", "", "path to user prompt template")
	outputTemplatePathArg := flag.String("output", "", "this flag should only be used when the prompt specifies JSON formatting.")
	flag.Parse()

	sumClient := genAi.NewGenAIClient(*genApiKindArg)
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

	if outputTemplatePathArg != nil && *outputTemplatePathArg != "" {
		if err := summarizer.SetOutputTemplate(*outputTemplatePathArg); err != nil {
			fmt.Println(summary)
			log.Fatal("Failed to format summary:", err)
		}
	}
	formattedSummary, err := summarizer.FormatOutput(summary)
	if err != nil {
		fmt.Println(summary)
		log.Fatal("Failed to format summary:", err)
	}
	fmt.Println(formattedSummary)
}
