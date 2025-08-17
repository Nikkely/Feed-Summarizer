package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var (
	// Command line flag for RSS feed URL
	fetchURL string
)

// FetchCmd represents the fetch command that retrieves RSS feeds as JSON
var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch RSS feed as JSON",
	Long: `Fetches an RSS feed from the specified URL and outputs it in JSON format.
The command supports standard RSS 2.0, RSS 1.0, and Atom formats.`,
	RunE: fetch,
}

func init() {
	// Register fetch command to root command
	rootCmd.AddCommand(FetchCmd)

	// Define command line flags
	FetchCmd.Flags().StringVarP(&fetchURL, "url", "u", "", "URL of the RSS feed to fetch")
	FetchCmd.MarkFlagRequired("url")
}

// fetch retrieves the RSS feed from the specified URL and outputs it as JSON.
// It handles the core functionality of the fetch command.
// Parameters:
//   - cmd: The Cobra command being run
//   - args: Command line arguments after the command name
//
// Returns:
//   - error: An error if the fetch operation fails, nil otherwise
func fetch(cmd *cobra.Command, args []string) error {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(fetchURL)
	if err != nil {
		return fmt.Errorf("failed to fetch RSS feed from URL %s: %w", fetchURL, err)
	}

	buf, err := json.Marshal(feed)
	if err != nil {
		return fmt.Errorf("failed to marshal feed to JSON: %w", err)
	}

	fmt.Println(string(buf))
	return nil
}
