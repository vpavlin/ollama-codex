package types

import "github.com/ollama/ollama/server"

type Manifest struct {
	server.Manifest
	Layers []Layer `json:"layers"`
	Config Layer   `json:"config"`
}

type Layer struct {
	server.Layer
	CID string `json:"cid"`
}
