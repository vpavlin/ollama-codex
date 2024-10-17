package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ollama/ollama/server"
	"github.com/spf13/cobra"
	"github.com/vpavlin/ollama-codex/types"
)

var manifestCid = ""

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pulls model from Codex",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		var manifest types.Manifest
		if manifestCid != "" {
			stream, err := download(manifestCid)
			if err != nil {
				log.Fatal(err)
			}

			data, err := io.ReadAll(stream)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(data, &manifest)
			if err != nil {
				log.Fatal(err)
			}
			mp := server.ParseModelPath(args[0])
			p, err := mp.GetManifestPath()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Storing manifest at", p)

			err = os.MkdirAll(filepath.Dir(p), 0o755)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(p, data, 0o644)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Downloading config", manifest.Config.Digest, "(CID:", manifest.Config.CID, ")")
		err := downloadLayer(manifest.Config)
		if err != nil {
			log.Fatal(err)
		}

		for _, layer := range manifest.Layers {
			log.Println("Downloading layer", layer.Digest, "(CID:", layer.CID, ")")

			err := downloadLayer(layer)
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	pullCmd.Flags().StringVar(&manifestCid, "manifest-cid", "", "Codex dataset CID of the model manifest")

	rootCmd.AddCommand(pullCmd)
}

func download(cid string) (io.ReadCloser, error) {
	url, err := url.JoinPath(codexApiUrl, "api/codex/v1/data", cid, "network")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func downloadLayer(layer types.Layer) error {
	stream, err := download(layer.CID)
	if err != nil {
		return err
	}
	defer stream.Close()

	filename, err := server.GetBlobsPath(layer.Digest)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(filename), 0o755)
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	c, err := io.Copy(fp, stream)
	if err != nil {
		return err
	}

	if c != layer.Size {
		return fmt.Errorf("size mismatch", c, "!=", layer.Size)
	}

	return nil
}
