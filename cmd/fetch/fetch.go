package main

import (
	"flag"
	"fmt"
	"log"
	rss "rss-summaries/internal"
)

func main() {
	feedUrl := flag.String("url", "https://example.com/feed", "rss feed url")
	flag.Parse()

	feedItems, err := rss.FetchRss(*feedUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(feedItems)
}
