package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// フラグ変数
	urlArg             string
	genAPIKindArg      string
	systemPromptPath   string
	userPromptPath     string
	formatOutput       bool
	outputTemplatePath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize RSS feed content using AI",
	Long: `A CLI tool that fetches RSS feed content and generates summaries using AI.
It supports custom prompts and output formatting.`,
	RunE: summarize,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// フラグの定義
	rootCmd.Flags().StringVarP(&urlArg, "url", "u", "", "RSS feed URL to summarize")
	rootCmd.Flags().StringVar(&genAPIKindArg, "gen-api-kind", "gemini", "Generative AI API type (currently only 'gemini' is supported)")
	rootCmd.Flags().StringVar(&systemPromptPath, "system-prompt", "", "Path to custom system prompt template file")
	rootCmd.Flags().StringVar(&userPromptPath, "user-prompt", "", "Path to custom user prompt template file")
	rootCmd.Flags().BoolVar(&formatOutput, "format", false, "Format output as JSON with template")
	rootCmd.Flags().StringVar(&outputTemplatePath, "output", "", "Custom output template path (only used when -format is true)")

	// 必須フラグの設定
	rootCmd.MarkFlagRequired("url")
}
