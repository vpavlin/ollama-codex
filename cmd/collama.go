package cmd

import "github.com/spf13/cobra"

var codexApiUrl = ""
var rootCmd = &cobra.Command{
	Use: "collama allows you to pull and push Ollama models from/to Codex",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&codexApiUrl, "codex-api-url", "a", "http://localhost:8080", "Codex API URL")
}

func Execute() error {
	return rootCmd.Execute()
}
