package main

import (
	"flag"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"sprayer/src/ui"
	"sprayer/src/ui/tui"
	"sprayer/src/version"
)

func main() {
	godotenv.Load()

	versionFlag := flag.Bool("version", false, "Print version information")
	shortVersionFlag := flag.Bool("v", false, "Print short version")
	tuiFlag := flag.Bool("tui", false, "Run in TUI mode")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("sprayer %s\n", version.WithPrefix())
		return
	}

	if *shortVersionFlag {
		fmt.Println(version.Version)
		return
	}

	if *tuiFlag {
		p := tea.NewProgram(tui.NewModel())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Run CLI mode
	runCLI()
}

func runCLI() {
	cli, err := ui.NewCLI()
	if err != nil {
		log.Fatal(err)
	}
	cli.Run()
}
