package main

import (
	"flag"
	"fmt"
	"log"
	sum "rss-summarizer/internal/summarize"
)

func main() {
	urlArg := flag.String("url", "https://example.com/feed", "rss feed url")
	genApiKindArg := flag.String("gen-api-kind", "gemini", "kind of generative API to use")
	flag.Parse()

	var sumClient sum.GenAIClient
	switch *genApiKindArg {
	case "gemini":
		sumClient = sum.NewGeminiClient("gemini-2.5-flash-lite")
	default:
		log.Fatal("Unsupported API")
	}

	summary, err := sum.Summarize(sumClient, *urlArg)
	if err != nil {
		log.Fatal("Failed to summarize URLs:", err)
	}

	fmt.Println(summary)
}
