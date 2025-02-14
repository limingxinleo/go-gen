package main

import (
	"embed"
	"github.com/limingxinleo/go-gen/cmd"
)

//go:embed .go-gen
var DefaultConfigDir embed.FS

func main() {
	cmd.Execute()
}
