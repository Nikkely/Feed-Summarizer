package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
)

// main is the entry point for the fetch command.
// It fetches an RSS feed from the provided URL and prints it as JSON.
//
// Flags:
//   - url: The URL of the RSS feed to fetch. Defaults to "https://example.com/feed".
//
// Example:
//   go run fetch.go -url="https://example.com/rss"
func main() {
	feedURL := flag.String("url", "https://example.com/feed", "rss feed url")
	flag.Parse()

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(*feedURL)
	if err != nil {
		log.Fatalf("failed to fetch RSS feed from URL %s: %v", *feedURL, err)
	}

	var buf []byte
	if buf, err = json.Marshal(feed); err != nil {
		log.Fatalf("failed to marshal feed: %v", err)
	}

	fmt.Println(string(buf))
}
