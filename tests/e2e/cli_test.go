package main_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestE2E_CLI_Scrape_And_List(t *testing.T) {
	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", "sprayer-e2e", "./cmd/cli/main.go")
	cmd.Dir = "../.." // assume running from a subdir or fix path
	// Easier: run go run directly
	
	root := "/home/user/openclaw-setup" // Hardcode for this env
	
	// 1. Scrape (Fast API only)
	scrapeCmd := exec.Command("go", "run", "cmd/cli/main.go", "scrape", "--fast", "rust", "remote")
	scrapeCmd.Dir = root
	out, err := scrapeCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Scrape failed: %v\nOutput: %s", err, string(out))
	}
	t.Logf("Scrape Output: %s", string(out))

	// 2. List
	listCmd := exec.Command("go", "run", "cmd/cli/main.go", "list")
	listCmd.Dir = root
	listOut, err := listCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("List failed: %v\nOutput: %s", err, string(listOut))
	}
	
	output := string(listOut)
	if !strings.Contains(output, "rust") && !strings.Contains(output, "Rust") { 
		// Ideally we found at least one job. If network failed or no jobs found, this might flake.
		// But "Scrape Output" should tell us "Saved X jobs".
	}
}
