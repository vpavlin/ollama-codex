package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ollama/ollama/server"
	"github.com/ollama/ollama/types/model"
	"github.com/spf13/cobra"
	"github.com/vpavlin/ollama-codex/types"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "pushes model to Codex",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		models, err := server.Manifests()
		if err != nil {
			log.Fatal(err)
		}

		//log.Println(models)
		name := model.ParseName(args[0])
		model, ok := models[name]
		if !ok {
			log.Fatal("unknown model ", args[0])
		}

		codexManifest := types.Manifest{Manifest: *model}

		log.Println("Found model", name, "with", len(model.Layers), "layers")
		for _, layer := range model.Layers {
			stream, err := layer.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer stream.Close()
			log.Printf("Starting to upload layer %s (size: %d)", layer.Digest, layer.Size)

			cid, err := upload(stream)

			log.Println(layer.Digest, " => ", string(cid))
			codexManifest.Layers = append(codexManifest.Layers, types.Layer{Layer: layer, CID: string(cid)})
		}

		stream, err := model.Config.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer stream.Close()

		cid, err := upload(stream)
		log.Println(model.Config.Digest, " => ", string(cid))
		codexManifest.Config.CID = cid
		codexManifest.Config.Digest = model.Config.Digest
		codexManifest.Config.From = model.Config.From
		codexManifest.Config.MediaType = model.Config.MediaType
		codexManifest.Config.Size = model.Config.Size

		//common.PrettyPrint(codexManifest)

		data, err := json.Marshal(codexManifest)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)

		cid, err = upload(reader)
		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile(strings.ReplaceAll(strings.ReplaceAll(name.String(), "/", "_"), ":", "_"), data, 0644)

		log.Println("Manifest CID", string(cid))
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}

func upload(stream io.Reader) (string, error) {
	url, err := url.JoinPath(codexApiUrl, "api/codex/v1/data")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, stream)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cid, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(cid), nil
}
