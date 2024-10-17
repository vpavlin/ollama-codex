package common

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/ollama/ollama/envconfig"
	"github.com/ollama/ollama/types/model"
)

const (
	defaultRegistryFilename = "registry.json"
)

type Registry struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name string `json:"name"`
	CID  string `json:"cid"`
}

func GetRawRegistryManifestURL(radicleApiUrl string, radID string, ref string, filename string) (string, error) {
	if filename == "" {
		filename = "manifest.json"
	}

	if ref == "" {
		ref = "head"
	}

	return url.JoinPath(radicleApiUrl, "raw", radID, ref, filename)
}

func NewModelRegistry(url string) (*Registry, error) {
	var registry Registry
	registryPath := path.Join(envconfig.Models(), defaultRegistryFilename)

	info, err := os.Stat(registryPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		} else {

		}
	}

	if info != nil && info.ModTime().After(time.Now().Add(-1*time.Minute)) {
		data, err := os.ReadFile(registryPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &registry)
		if err != nil {
			return nil, err
		}

		log.Println("Using cached registry manifest")

		return &registry, nil
	}

	resp, err := http.Get(url)
	if err != nil {

	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch registry manifest: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(registryPath, data, 0644)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &registry)
	if err != nil {
		return nil, err
	}

	return &registry, nil
}

func (r *Registry) GetCID(name string) (string, error) {
	fullName := model.ParseName(name)
	for _, m := range r.Models {
		if m.Name == fullName.String() {
			return m.CID, nil
		}
	}

	return "", fmt.Errorf("model not found in registry %s", fullName)
}
