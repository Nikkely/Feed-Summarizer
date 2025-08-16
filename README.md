# RSS Summarizer

RSS Summarizer is a Go-based application designed to fetch RSS feeds and generate summaries. This project leverages AI to efficiently summarize information.

## Key Features
- Fetching RSS feeds
- Fetching HTML pages
- Generating summaries (using AI)
  - Currently supports only the Gemini model, but plans to support multiple models in the future.
- Generating coverage reports and deploying them to GitHub Pages

## Directory Structure
```
cmd/
  fetch/         # Command for fetching RSS feeds
  summarize/     # Command for generating summaries
docs/            # Documentation
internal/
  ai_client/     # Implementation of the AI client
  fetcher/       # Logic for fetching feeds and HTML pages
  summarize/     # Logic for generating summaries
pkg/
  prompt/        # Logic for generating prompts
scripts/         # Scripts for building and testing
templates/       # Prompt templates
```

## Required Environment Variables
- `GEMINI_API_KEY`: API key for using Gemini for summary generation.

## Setup
1. Install the required dependencies:
   ```sh
   go mod tidy
   ```
2. Set the required environment variables:
   ```sh
   export GEMINI_API_KEY=your_api_key_here
   ```

## Usage
### Generating Summaries
You can generate summaries using the following command:
```sh
go run cmd/summarize/summarize.go --url <feed_url> 
```

## License
This project is licensed under the MIT License.
