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
	systemPromptArg := flag.String("system-prompt", "templates/system_prompt.txt", "path to system prompt template")
	userPromptArg := flag.String("user-prompt", "templates/user_prompt.tmpl", "path to user prompt template")
	flag.Parse()

	var sumClient genAi.GenAIClient
	switch *genApiKindArg {
	case "gemini":
		sumClient = genAi.NewGeminiClient("gemini-2.5-flash-lite")
	default:
		log.Fatal("Unsupported API")
	}

	summarizer := sum.NewSummarizer(sumClient, fetcher.FetchFeed, fetcher.FetchHTML)
	if err := summarizer.LoadPromptBuilder(*systemPromptArg, *userPromptArg); err != nil {
		log.Fatal("Failed to load prompt builder:", err)
	}
	summary, err := summarizer.Summarize(*urlArg)
	if err != nil {
		log.Fatal("Failed to summarize URLs:", err)
	}

	fmt.Println(summary)
}
