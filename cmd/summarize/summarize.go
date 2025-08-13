package main

import (
	"flag"
	"fmt"
	"log"
	sum "rss-summarizer/internal/summarize"
	"strings"
)

func main() {
	urlsArg := flag.String("urls", "https://example.com/feed", "rss feed url")
	separatorArg := flag.String("sep", ",", "separator for multiple urls")
	genApiKindArg := flag.String("gen-api-kind", "gemini", "kind of generative API to use")
	flag.Parse()

	var sumClient sum.GenAIClient
	switch *genApiKindArg {
	case "gemini":
		sumClient = sum.NewGeminiClient("gemini-2.5-flash-lite")
	default:
		log.Fatal("Unsupported API")
	}

	summary, err := sum.Summarize(sumClient, strings.Split(*urlsArg, *separatorArg))
	if err != nil {
		log.Fatal("Failed to summarize URLs:", err)
	}

	fmt.Println(summary)
}
