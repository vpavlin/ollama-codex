package main

import (
	"log"

	"github.com/vpavlin/ollama-codex/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
