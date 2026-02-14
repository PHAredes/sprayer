package main

import (
	"log"

	"sprayer/internal/ui"
)

func main() {
	cli, err := ui.NewCLI()
	if err != nil {
		log.Fatal(err)
	}
	cli.Run()
}