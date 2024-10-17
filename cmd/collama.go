package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vpavlin/ollama-codex/common"
)

var (
	codexApiUrl   = ""
	radicleApiUrl = ""
	rid           = ""
	ref           = ""
	filename      = ""
)

var registry *common.Registry

var rootCmd = &cobra.Command{
	Use: "collama allows you to pull and push Ollama models from/to Codex",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&codexApiUrl, "codex-api-url", "a", "http://localhost:8080", "Codex API URL")
	rootCmd.PersistentFlags().StringVar(&radicleApiUrl, "radicle-api-url", "https://seed.radicle.garden/", "Radicle Seed Node API URL (to fetch model name -> CID mappings)")
	rootCmd.PersistentFlags().StringVar(&rid, "radicle-id", "rad:z4UyZVeb2hy7oKxZ1mUgR1ujAzFsy", "Radicle repository ID")
	rootCmd.PersistentFlags().StringVar(&ref, "radicle-ref", "head", "Radicle repository ref")
	rootCmd.PersistentFlags().StringVar(&filename, "radicle-filename", "manifest.json", "Name of the file containing the registry data")

}

func Execute() error {
	manifestUrl, err := common.GetRawRegistryManifestURL(radicleApiUrl, rid, ref, filename)
	if err != nil {
		log.Fatal(err)
	}
	registry, err = common.NewModelRegistry(manifestUrl)
	if err != nil {
		log.Fatal(err)
	}
	return rootCmd.Execute()
}
