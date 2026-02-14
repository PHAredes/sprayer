package main

import (
	"log"

	"github.com/joho/godotenv"
	"sprayer/internal/ui"
)

func main() {
	godotenv.Load()
	cli, err := ui.NewCLI()
	if err != nil {
		log.Fatal(err)
	}
	cli.Run()
}