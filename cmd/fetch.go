package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

// FetchCmd represents the fetch command that retrieves RSS feeds as JSON
var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Only Fetch RSS feed as JSON",
	Long: `Fetches an RSS feed from the specified URL and outputs it in JSON format.
The command supports standard RSS 2.0, RSS 1.0, and Atom formats.`,
	Args: cobra.MinimumNArgs(1),
	RunE: fetch,
}

func init() {
	// Register fetch command to root command
	rootCmd.AddCommand(FetchCmd)
}

// fetch retrieves the RSS feed from the specified URL and outputs it as JSON.
// It handles the core functionality of the fetch command.
// Parameters:
//   - _: The Cobra command being run
//   - args: Command line arguments after the command name
//
// Returns:
//   - error: An error if the fetch operation fails, nil otherwise
func fetch(_ *cobra.Command, args []string) error {
	fp := gofeed.NewParser()
	if len(args) < 1 {
		return fmt.Errorf("missing URL argument")
	}

	for _, url := range args {
		feed, err := fp.ParseURL(url)
		if err != nil {
			return fmt.Errorf("failed to fetch feed from URL %s: %w", url, err)
		}

		buf, err := json.Marshal(feed)
		if err != nil {
			return fmt.Errorf("failed to marshal feed to JSON: %w", err)
		}

		fmt.Println(string(buf))
	}
	return nil
}
